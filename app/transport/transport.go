package transport

import (
	"github.com/gorilla/mux"
	"github.com/jackvonhouse/car-enrichment/app/usecase"
	_ "github.com/jackvonhouse/car-enrichment/docs"
	"github.com/jackvonhouse/car-enrichment/internal/transport/car"
	"github.com/jackvonhouse/car-enrichment/internal/transport/router"
	"github.com/jackvonhouse/car-enrichment/pkg/log"
	"github.com/swaggo/http-swagger/v2"
)

type Transport struct {
	router *router.Router
}

func New(
	useCase usecase.UseCase,
	logger log.Logger,
) Transport {

	transportLogger := logger.WithField("layer", "transport")

	r := router.New("/api/v1")

	r.Handle(map[string]router.Handlify{
		"/car": car.New(useCase.Car, transportLogger),
	})

	r.Router().
		PathPrefix("/swagger").
		Handler(httpSwagger.WrapHandler)

	return Transport{
		router: r,
	}
}

func (t Transport) Router() *mux.Router { return t.router.Router() }
