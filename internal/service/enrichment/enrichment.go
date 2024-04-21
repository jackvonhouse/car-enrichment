package enrichment

import (
	"context"
	"encoding/json"
	"github.com/jackvonhouse/car-enrichment/config"
	"github.com/jackvonhouse/car-enrichment/internal/dto"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
	"net/http"
	"sync"
	"time"
)

type Service struct {
	config config.API

	logger log.Logger
}

func New(
	config config.API,
	logger log.Logger,
) Service {

	return Service{
		config: config,
		logger: logger.WithField("unit", "enrichment"),
	}
}

func (e Service) Enrichment(
	_ context.Context,
	regNumbers []string,
) (map[int64]dto.Car, error) {

	const maxAttempts = 3

	e.logger.Debugf("starting enrichment for %d cars", len(regNumbers))
	e.logger.Debugf("max attempts: %d", maxAttempts)

	cars := make(map[int64]dto.Car)

	m := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(len(regNumbers))

	for i, regNumber := range regNumbers {
		e.logger.Debugf("starting enrichment for car with regNum %s", regNumber)

		i := int64(i)

		go func(i int64, regNumber string) {
			defer wg.Done()

			var (
				car dto.Car
				err error
			)

			for attempt := 0; attempt < maxAttempts; attempt++ {
				car, err = e.makeRequest(regNumber)
				if err == nil {
					m.Lock()
					cars[i] = car
					m.Unlock()

					break
				}
			}
		}(i, regNumber)
	}

	wg.Wait()

	return cars, nil
}

func (e Service) makeRequest(
	regNumber string,
) (dto.Car, error) {

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	request, err := http.NewRequest(http.MethodGet, e.config.Url, nil)
	if err != nil {
		return dto.Car{}, err
	}

	queries := request.URL.Query()
	queries.Add("regNum", regNumber)
	request.URL.RawQuery = queries.Encode()

	response, err := client.Do(request)
	if err != nil {
		return dto.Car{}, err
	}

	defer response.Body.Close()

	car := dto.Car{}

	if err := json.NewDecoder(response.Body).Decode(&car); err != nil {
		return dto.Car{}, err
	}

	return car, nil
}
