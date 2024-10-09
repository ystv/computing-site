package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ystv/computing_site/link"
	"github.com/ystv/computing_site/team"
	"github.com/ystv/computing_site/templates"
)

type (
	Router struct {
		address   string
		templater *templates.Templater
		link      *link.Link
		team      *[]team.Member
		commit    string
		version   string
		cert      []byte
		key       []byte
		router    *echo.Echo
	}

	RouterConf struct {
		Address string
		Link    *link.Link
		Team    *[]team.Member
		Commit  string
		Version string
		cert    []byte
		key     []byte
	}
)

func NewRouter(conf *RouterConf) *Router {
	r := &Router{
		address: conf.Address,
		link:    conf.Link,
		team:    conf.Team,
		commit:  conf.Commit,
		version: conf.Version,
		router:  echo.New(),
	}
	r.router.HideBanner = true

	r.middleware()

	r.loadRoutes()

	return r
}

func (r *Router) Start() error {
	r.router.Logger.Error(r.router.StartTLS(r.address, r.cert, r.key))
	return fmt.Errorf("failed to start router on address %s", r.address)
}

// middleware initialises web server middleware
func (r *Router) middleware() {
	r.router.Pre(middleware.RemoveTrailingSlash())
	r.router.Use(middleware.Recover())
	r.router.Use(middleware.BodyLimit("15M"))
	r.router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
}

func (r *Router) loadRoutes() {
	r.router.RouteNotFound("/*", func(c echo.Context) error { return c.Redirect(http.StatusFound, "/") })

	assetHandler := http.FileServer(http.FS(echo.MustSubFS(embeddedFiles, "public")))

	r.router.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", assetHandler)))

	base := r.router.Group("/")

	base.GET("api/health", func(c echo.Context) error {
		marshal, err := json.Marshal(struct {
			Status int `json:"status"`
		}{
			Status: http.StatusOK,
		})
		if err != nil {
			log.Println(err)
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  err.Error(),
				Internal: err,
			}
		}

		c.Response().Header().Set("Content-Type", "application/json")
		return c.JSON(http.StatusOK, marshal)
	})

	base.GET("", r.HomeFunc)
}

func (r *Router) HomeFunc(c echo.Context) error {
	params := &templates.DashboardParams{
		Link:    r.link,
		Team:    r.team,
		Commit:  r.commit,
		Version: r.version,
	}

	return r.templater.RenderTemplate(c.Response(), params, "dashboard.tmpl")
}
