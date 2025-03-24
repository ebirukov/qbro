package broker

import (
	"context"
	httpApi "ebirukov/qbro/internal/api/http"
	"ebirukov/qbro/internal/config"
	"ebirukov/qbro/internal/connector"
	"ebirukov/qbro/internal/service"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type BrokerApp struct {
	httpServer *http.Server
	brokerSvc  *service.BrokerSvc
	qRegistry  *service.QConnRegistry

	appCtx       context.Context
	appCtxCancel context.CancelCauseFunc
}

func NewBrokerApp(cfg config.BrokerCfg) (*BrokerApp, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid app configuration: %w", err)
	}

	app := &BrokerApp{}

	app.appCtx, app.appCtxCancel = context.WithCancelCause(context.Background())

	app.qRegistry = service.NewQueueConnRegistry(cfg, connector.NewChanQueueConnector(cfg.MaxQueueSize))

	app.brokerSvc = service.NewBrokerSvc(app.appCtx, cfg, app.qRegistry)

	httpHandlers := httpApi.NewBrokerHandler(cfg, app.brokerSvc)
	router := httpApi.RegisterRoutes(httpHandlers)

	app.httpServer = &http.Server{
		Addr:    cfg.HttpAddr(),
		Handler: router,
	}

	return app, nil
}

func (app *BrokerApp) Run() error {
	closeSig := make(chan os.Signal, 1)
	defer close(closeSig)

	errSig := make(chan error)
	defer close(errSig)

	go func() {
		err := app.httpServer.ListenAndServe()
		if err != nil {
			switch {
			case errors.Is(err, http.ErrServerClosed):
			default:
				errSig <- err
				return
			}
		}

		errSig <- nil
	}()

	signal.Notify(closeSig, syscall.SIGINT, syscall.SIGTERM)

	var err error

	select {
	case <-closeSig:
	case err = <-errSig:
	}

	// cancel all requests that read from queues
	app.appCtxCancel(errors.New("app was stop"))

	return errors.Join(
		err,
		app.qRegistry.Shutdown(),
		app.httpServer.Shutdown(context.Background()),
	)
}
