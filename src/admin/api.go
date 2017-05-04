package admin

func (server *AdminServer) initAPIRoute() {

    server.e.GET("/api/apis", server.getAPIs())
}
