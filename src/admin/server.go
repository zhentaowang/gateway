package admin

import (
    "github.com/labstack/echo"
    "fmt"
    mw "github.com/labstack/echo/middleware"
    "gateway/src/model"
)

type Result struct {
    Code  int         `json:"code, omitempty"`
    Value interface{} `json:"value"`
}

type AdminServer struct {
    e    *echo.Echo
    user string
    pwd  string
    address string
    store model.Store
}

func NewAdminServer(address string, user string, pwd string, store model.Store) *AdminServer {
    adminServer := &AdminServer{
        user: user,
        pwd:  pwd,
        e:    echo.New(),
        address: address,
        store: store,
    }

    adminServer.initHTTPServer()

    return adminServer
}

func (server *AdminServer) initHTTPServer() {
    server.e.Use(mw.Logger())
    server.e.Use(mw.Recover())
    server.e.Use(mw.Gzip())
    server.e.Use(mw.BasicAuth(func(inputUser string, inputPwd string, c echo.Context) (error, bool) {
        if inputUser == server.user && server.pwd == inputPwd {
            return nil, true
        }
        return nil, false
    }))

    server.e.Static("/assets", "admin/public/assets")
    server.e.Static("/html", "admin/public/html")

    server.e.File("/", "admin/public/html/base.html")

    server.initAPIRoute()
}

func (server *AdminServer) Start() {
    fmt.Printf("start at %s\n", server.address)
    server.e.Logger.Fatal(server.e.Start(server.address))
}
