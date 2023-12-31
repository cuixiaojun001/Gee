package main

/*
   (1) $ curl -i http://localhost:9999/
   (2) $ curl "http://localhost:9999/hello?name=geektutu"
   (3) $ curl "http://localhost:9999/login" -X POST -d 'username=geektutu&password=1234'
   (4) $ curl "http://localhost:9999/xxx"
*/

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"cxj"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":8080")
}
