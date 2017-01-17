package loket

import (
	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"
	"github.com/mataharimall/micro"
	"net/http"
)

type LoketRoute struct{}

func init() {
	micro.RouterManager.Register("route.loket", &LoketRoute{})
}

func (l *LoketRoute) SetRoute(e *echo.Echo) *echo.Echo {
	e.Use(em.Logger())
	e.Post("/loket/event", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	return e
}
