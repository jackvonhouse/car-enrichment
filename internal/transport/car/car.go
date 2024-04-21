package car

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jackvonhouse/car-enrichment/internal/dto"
	"github.com/jackvonhouse/car-enrichment/internal/transport"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
	"net/http"
	"strings"
	"time"
)

type carUseCase interface {
	Create(context.Context, dto.CreateCar) (map[int64]string, error)

	Get(context.Context, dto.Filter, dto.Pagination) ([]dto.Car, error)

	Update(context.Context, dto.Car) error

	Delete(context.Context, int64) error
}

type Transport struct {
	car carUseCase

	logger log.Logger
}

func New(
	car carUseCase,
	logger log.Logger,
) Transport {
	return Transport{
		car:    car,
		logger: logger.WithField("unit", "car"),
	}
}

func (t Transport) loggerMiddleware(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.logger.Infof("[%s] received request %s", r.Method, r.URL.String())

		next.ServeHTTP(w, r)
	})
}

func (t Transport) Handle(
	router *mux.Router,
) {
	logRouter := router.PathPrefix("").Subrouter()
	logRouter.Use(t.loggerMiddleware)

	logRouter.HandleFunc("", t.Create).
		Methods(http.MethodPost)

	logRouter.HandleFunc("", t.Get).
		Methods(http.MethodGet)

	logRouter.HandleFunc("/{id:[0-9]+}", t.Update).
		Methods(http.MethodPut)

	logRouter.HandleFunc("/{id:[0-9]+}", t.Delete).
		Methods(http.MethodDelete)
}

// Create godoc
// @Summary			Создание автомобиля
// @Description		Создание и обогощение автомобиля
// @Accept			json
// @Produce			json
// @Param			request body dto.CreateCar true "Массив гос. номеров"
// @Success			200 {object} object{result=bool}
// @Failure			409 {object} object{error=string} "Автомобиль или владелец уже существует"
// @Failure			500 {object} object{error=string} "Неизвестная ошибка"
// @Tags			Автомобиль
// @Router /car [post]
func (t Transport) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	data := dto.CreateCar{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		transport.Error(w, http.StatusInternalServerError, "invalid json structure")

		return
	}

	if len(data.RegNumbers) == 0 {
		transport.Error(w, http.StatusBadRequest, "empty registration numbers")

		return
	}

	for _, regNumber := range data.RegNumbers {
		if len(strings.TrimSpace(regNumber)) == 0 {
			transport.Error(w, http.StatusBadRequest, "empty registration number")

			return
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	failed, err := t.car.Create(ctx, data)
	if len(failed) != 0 {
		transport.Response(w, map[string]any{
			"error":   "some cars are not enriched",
			"regNums": failed,
		})

		return
	}

	if err != nil {
		t.logger.Warn(err)

		code, msg := transport.ErrorToHttpResponse(
			err,
			transport.DefaultErrorHttpCodes,
		)

		transport.Error(w, code, msg)

		return
	}

	transport.Response(w, map[string]any{"success": true})
}

// Get godoc
// @Summary			Получить автомобилей
// @Description		Получение автомобилей с возможностью фильтрации
// @Accept			json
// @Produce			json
// @Param			limit query int false "Лимит"
// @Param			offset query int false "Смещение"
// @Param			regNum query string false "Гос. номер"
// @Param			mark query string false "Марка"
// @Param			model query string false "Модель"
// @Param			year query int false "Год"
// @Param			ownerName query string false "Имя владельца"
// @Param			ownerSurname query string false "Фамилия владельца"
// @Param			ownerPatronymic query string false "Отчество владельца"
// @Success			200 {array} dto.Car
// @Failure			404 {object} object{error=string} "Автомобили отсутствуют"
// @Failure			500 {object} object{error=string} "Неизвестная ошибка"
// @Tags			Автомобиль
// @Router /car [get]
func (t Transport) Get(
	w http.ResponseWriter,
	r *http.Request,
) {

	queries := r.URL.Query()

	regNum := queries.Get("regNum")
	mark := queries.Get("mark")
	model := queries.Get("model")
	ownerName := queries.Get("ownerName")
	ownerSurname := queries.Get("ownerSurname")
	ownerPatronymic := queries.Get("ownerPatronymic")

	year, err := transport.StringToInt(queries.Get("year"))
	if err != nil || year < 0 {
		year = 0
	}

	filter := dto.Filter{
		RegNum:          regNum,
		Mark:            mark,
		Model:           model,
		Year:            year,
		OwnerName:       ownerName,
		OwnerSurname:    ownerSurname,
		OwnerPatronymic: ownerPatronymic,
	}

	limit, err := transport.StringToInt(queries.Get("limit"))
	if err != nil || limit <= 0 {
		// В зависимости от логики выбрасывать ошибку
		// или устанавливать limit по умолчанию
		// transport.Error(w, http.StatusBadRequest, "invalid limit")
		// return

		limit = 10
	}

	offset, err := transport.StringToInt(queries.Get("offset"))
	if err != nil || offset < 0 {
		// В зависимости от логики выбрасывать ошибку
		// или устанавливать offset по умолчанию
		// transport.Error(w, http.StatusBadRequest, "invalid offset")
		// return

		offset = 0
	}

	pagination := dto.Pagination{
		Limit:  limit,
		Offset: offset,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	cars, err := t.car.Get(ctx, filter, pagination)
	if err != nil {
		t.logger.Warn(err)

		code, msg := transport.ErrorToHttpResponse(
			err,
			transport.DefaultErrorHttpCodes,
		)

		transport.Error(w, code, msg)

		return
	}

	transport.Response(w, cars)
}

// Update godoc
// @Summary			Обновить автомобиль
// @Description		Обновление автомобиля
// @Accept			json
// @Produce			json
// @Param			request body dto.Car true "Данные об автомобиле"
// @Param			id path int true "Идентификатор автомобиля"
// @Success			200 {object} object{result=bool}
// @Failure			404 {object} object{error=string} "Автомобиль или владелец не найдены"
// @Failure			409 {object} object{error=string} "Автомобиль уже существует"
// @Failure			500 {object} object{error=string} "Неизвестная ошибка"
// @Tags			Автомобиль
// @Router /car/{id} [put]
func (t Transport) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	vars := mux.Vars(r)

	carId, err := transport.StringToInt(vars["id"])
	if err != nil || carId <= 0 {
		transport.Error(w, http.StatusBadRequest, "invalid car id")

		return
	}

	data := dto.Car{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		transport.Error(w, http.StatusInternalServerError, "invalid json structure")

		return
	}

	data.ID = int64(carId)

	if data.Year != 0 && (data.Year < 1900 || data.Year > time.Now().Year()) {
		transport.Error(w, http.StatusBadRequest, "invalid year")

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := t.car.Update(ctx, data); err != nil {
		t.logger.Warn(err)

		code, msg := transport.ErrorToHttpResponse(
			err,
			transport.DefaultErrorHttpCodes,
		)

		transport.Error(w, code, msg)

		return
	}

	transport.Response(w, map[string]any{"success": true})
}

// Delete godoc
// @Summary			Удалить автомобиль
// @Description		Удаление автомобиля
// @Accept			json
// @Produce			json
// @Param			id path int true "Идентификатор автомобиля"
// @Success			200 {object} object{result=bool}
// @Failure			404 {object} object{error=string} "Автомобиль не найден"
// @Failure			500 {object} object{error=string} "Неизвестная ошибка"
// @Tags			Автомобиль
// @Router /car/{id} [delete]
func (t Transport) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	vars := mux.Vars(r)

	carId, err := transport.StringToInt(vars["id"])
	if err != nil || carId <= 0 {
		transport.Error(w, http.StatusBadRequest, "invalid car id")

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := t.car.Delete(ctx, int64(carId)); err != nil {
		t.logger.Warn(err)

		code, msg := transport.ErrorToHttpResponse(
			err,
			transport.DefaultErrorHttpCodes,
		)

		transport.Error(w, code, msg)

		return
	}

	transport.Response(w, map[string]any{"success": true})
}
