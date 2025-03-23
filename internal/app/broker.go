package broker

import (
	"context"
	"ebirukov/qbro/internal/adapter"
	httpApi "ebirukov/qbro/internal/api/http"
	"ebirukov/qbro/internal/model"
	"ebirukov/qbro/internal/service"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type BrokerApp struct {
	httpServer *http.Server
	brokerSvc *service.BrokerSvc
	qRegistry *service.QRegistry

	appCtx context.Context
	appCtxCancel context.CancelCauseFunc
}

func NewBrokerApp(cfg model.Config) *BrokerApp {
	app := &BrokerApp{}

	app.appCtx, app.appCtxCancel = context.WithCancelCause(context.Background())

	app.qRegistry = service.NewQueueRegistry(app.appCtx, cfg, adapter.NewChanQueueCreator())

	app.brokerSvc = service.NewBrokerSvc(cfg, app.qRegistry)

	httpHandlers := httpApi.NewBrokerHandler(cfg, app.brokerSvc)
	router := httpApi.RegisterRoutes(httpHandlers)

	port := cfg.HttpPort
	if port == "" {
		port = ":8080"
	}

	app.httpServer = &http.Server{
		Addr: port,
		Handler: router,
	}

	return app
}

func (app *BrokerApp) Run() {
	go func() {
		err := app.httpServer.ListenAndServe()
		if err != nil {
			switch {
			case errors.Is(err, http.ErrServerClosed):
			default:
				log.Fatal(err)
			}
		}
	}()

	closeSig := make(chan os.Signal, 1)
	signal.Notify(closeSig, syscall.SIGINT, syscall.SIGTERM)

	<- closeSig

	app.appCtxCancel(errors.New("app was stop"))

	app.qRegistry.Shutdown()
	app.httpServer.Shutdown(app.appCtx)
}