package main

import (
	"net/http"
)

func main() {

	connectGRPC()
	defer grpcConn.Close()

	mux := routes()

	go actor()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	logger.Info(
		"Server started",
		"address", ":8080",
	)

	go func() {

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			logger.Error(
				"Server failed",
				"error", err,
			)
		}
	}()

	waitForShutdown(server)
}
