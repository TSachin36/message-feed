package api

func Run() {

	Start()

	waitForShutdown(httpServer)

	if grpcConn != nil {
		grpcConn.Close()
	}
}
