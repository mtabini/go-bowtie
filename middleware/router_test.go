package middleware

import (
	"github.com/mtabini/go-bowtie"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	r := NewRouter()

	r.GET("/:id", func(c bowtie.Context) {
		c.Response().Write([]byte("Hello " + c.(*RouterContext).Params.ByName("id")))
	})

	s := bowtie.NewServer()

	s.AddMiddlewareProvider(r)

	ss := httptest.NewServer(s)
	defer ss.Close()

	res, err := http.Get(ss.URL + "/test")

	if err != nil {
		t.Errorf("Unable to run test server: %s", err)
	}

	defer res.Body.Close()

	output, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("Unable to run test server: %s", err)
	}

	if string(output) != "Hello test" {
		t.Errorf("Unexpected response from test server: %s", output)
	}
}
