package api

import (
	"log"
	"net/http"

	_ "net/http/pprof"
)

func startPprof() {

	go func() {

		log.Println("pprof running on :6060")

		err := http.ListenAndServe(
			"localhost:6060",
			nil,
		)

		if err != nil {
			log.Printf(
				"pprof server stopped: %v",
				err,
			)
		}
	}()
}
