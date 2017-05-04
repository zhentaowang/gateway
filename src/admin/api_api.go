package admin

import (
    "github.com/labstack/echo"
    "net/http"
    "github.com/labstack/gommon/log"
)

func (server *AdminServer) getAPIs() echo.HandlerFunc {
    return func(c echo.Context) error {
        values, err := server.store.GetAPIs()
        if err != nil {
            log.Fatal(err)
        }

        return c.JSON(http.StatusOK, &Result{
            Code: 0,
            Value: values,
        })
    }
}
