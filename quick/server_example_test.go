package quick

import (
	"github.com/mtabini/go-bowtie"
	"github.com/mtabini/go-bowtie/middleware"
	"net/http"
)

func Example_quick_server() {
	// Create a new Bowtie quick server
	s := New()

	s.GET("/test/:id", func(c bowtie.Context) {
		id := c.(*middleware.RouterContext).Params.ByName("id")

		c.Response().WriteString("The ID is " + id)
	})

	http.ListenAndServe(":8000", s)
}
