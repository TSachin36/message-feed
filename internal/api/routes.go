package api

import "net/http"

func routes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.Handle(
		"/messages",
		loggingMiddleware(http.HandlerFunc(messagesHandler)),
	)

	mux.Handle(
		"/about/",
		loggingMiddleware(
			http.StripPrefix(
				"/about/",
				http.FileServer(http.Dir("web/about")),
			),
		),
	)

	mux.Handle(
		"/list",
		loggingMiddleware(
			http.HandlerFunc(listHandler),
		),
	)

	mux.Handle(
		"/ws",
		loggingMiddleware(
			http.HandlerFunc(websocketHandler),
		),
	)

	return mux
}
