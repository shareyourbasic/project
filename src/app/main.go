package main

import (
	"io"
	"net/http"
	"html/template"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/thoas/stats"
	"fmt"
	"time"
)

var (
	globalData *Data
)

type (
	Template struct {
		templates *template.Template
	}
)

func (t *Template) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.Execute(w, data)
}

func update(c *echo.Context) error {
	globalData.UpdateDB(true)
	globalData.Calculate()
	return c.JSONIndent(200, globalData, "", "  ")
}

func raw(c *echo.Context) error {
	globalData.Calculate()
	return c.JSONIndent(200, globalData, "", "  ")
}

func welcome(c *echo.Context) error {
	dat := globalData.GetTemplate()
	fmt.Println("template: %v", dat)
	return c.Render(http.StatusOK, "welcome", dat)
}

func main() {
	var err error
	globalData, err = LoadData()
	if err != nil {
		fmt.Println("data.json not found")
		panic(err)
	}
	go globalData.UpdateDB(false)

	e := echo.New()

	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(mw.Gzip())

	// https://github.com/thoas/stats
	s := stats.New()
	e.Use(s.Handler)
	// Route
	e.Get("/stats", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, s.Data())
	})

	e.Static("/css", "public/css")
	e.Static("/fonts", "public/fonts")
	e.Static("/img", "public/img")
	e.Static("/js", "public/js")

	t := &Template{
		templates: template.Must(template.ParseFiles("public/index.html")),
	}
	e.SetRenderer(t)
	e.Get("/update", update)
	e.Get("/raw", raw)
	e.Get("/", welcome)
	fmt.Println("starting on :80")
	go updateDB()
	// Start server
	e.Run(":80")
}

func updateDB() {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	count := 0
	for {
		select {
		case <-ticker.C:
			count++
			if count % 10 == 0 {
				go globalData.UpdateDB(true)
			}else {
				go globalData.UpdateDB(false)
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}