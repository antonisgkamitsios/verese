package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *Application) Serve(routes http.Handler) error {
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      routes,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutDownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// block while waiting for notify
		s := <-quit

		app.Logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutDownError <- err
		}

		// close background tasks if any
		bgTasksClosed := make(chan struct{})

		go func() {
			app.wg.Wait()
			bgTasksClosed <- struct{}{}
		}()

		select {
		case <-bgTasksClosed:
			app.Logger.Info("background tasks completed successfully", "addr", server.Addr)
			shutDownError <- nil
		case <-time.After(10 * time.Second):
			app.Logger.Error("background tasks are taking too much to complete, shutting down server...", "addr", server.Addr)
			shutDownError <- fmt.Errorf("timeout for closing background tasks exceeded")
		}

	}()

	err := server.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// block while we wait for an error to come back
	err = <-shutDownError
	if err != nil {
		return err
	}

	app.Logger.Info("stopped server", "addr", server.Addr)

	return nil
}
