package api

import (
	"net/http"
)

var httpServer *http.Server

func Start() {

	connectGRPC()

	startPprof()

	mux := routes()

	go actor()

	httpServer = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	logger.Info(
		"Server started",
		"address", ":8080",
	)

	go func() {

		if err := httpServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			logger.Error(
				"Server failed",
				"error", err,
			)
		}
	}()
}
