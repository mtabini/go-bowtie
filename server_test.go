package bowtie

import (
	"net/http"
	"testing"
	"time"
)

type testContext struct {
	Context
	t1 time.Time
	t2 time.Time
	t3 time.Time
	t4 time.Time
}

var tm = false
var tnm = false

func testMiddleware(c Context, next func()) {
	cc := c.(*testContext)
	cc.t4 = time.Now()
	tm = true

	cc.t1 = time.Now()
}

func testNextMiddleware(c Context, next func()) {
	cc := c.(*testContext)
	cc.t3 = time.Now()

	next()

	cc.t2 = time.Now()

	tnm = true
}

func TestServer(t *testing.T) {
	s := NewServer()

	s.AddMiddleware(testNextMiddleware)
	s.AddMiddleware(testMiddleware)

	s.AddContextFactory(func(c Context) Context {
		return &testContext{
			Context: c,
		}
	})

	r := &http.Request{}
	w := newMockWriter()

	s.ServeHTTP(w, r)

	if !tm {
		t.Error("The first middleware didn't run")
	}

	if !tnm {
		t.Error("The second middleware didn't run")
	}

	c := s.NewContext(r, w)

	s.Run(c)

	cc := c.(*testContext)

	if cc.t1.After(cc.t2) {
		t.Error("The next() handler doesn't seem to work")
	}

	if cc.t3.After(cc.t4) {
		t.Error("Middlewares doen't seem to be run in the proper order")
	}
}
