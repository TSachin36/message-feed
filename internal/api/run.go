package api

func Run() {

	Start()

	waitForShutdown(httpServer)

	closeGRPC()
}
