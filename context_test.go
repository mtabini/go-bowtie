package bowtie

import (
	"errors"
	"net/http"
	"testing"
)

type localContext struct {
	Context
	x string
}

func newLocalContext(r *http.Request, w http.ResponseWriter) *localContext {
	return &localContext{
		Context: NewContext(r, w),
		x:       "",
	}
}

type mockWriter struct {
	header  http.Header
	written []byte
	status  int
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		header:  http.Header{},
		written: []byte{},
		status:  0,
	}
}

func (m *mockWriter) Header() http.Header {
	return m.header
}

func (m *mockWriter) Write(p []byte) (int, error) {
	m.written = append(m.written, p...)
	return len(p), nil
}

func (m *mockWriter) WriteHeader(status int) {
	m.status = status
}

func TestContext(t *testing.T) {
	r := &http.Request{}
	w := newMockWriter()
	c := newLocalContext(r, w)

	if c.Request() != r {
		t.Error("Unexpectedly received the wrong request object")
	}

	if c.Response().Written() {
		t.Error("Context unexpectedly reports that it has data")
	}

	n, err := c.Response().Write([]byte("abc"))

	if err != nil {
		t.Errorf("Unable to write to context: %s", err)
	}

	if n != 3 {
		t.Errorf("Expected 3 bytes written, got %d instead", n)
	}

	w.written = []byte{}

	n, err = c.Response().WriteJSON(map[string]interface{}{"test": 123})

	if err != nil {
		t.Errorf("Unable to write JSON to context: %s", err)
	}

	if n != 12 {
		t.Errorf("Expected 12 bytes written, got %d instead", n)
	}

	if len(c.Response().Errors()) > 0 {
		t.Errorf("Context unexpectedly has errors after writing JSON: %#v", c.Response().Errors)
	}

	c.Response().WriteJSONOrError(map[string]interface{}{"test": 123}, errors.New("Error"))

	if len(c.Response().Errors()) == 0 {
		t.Error("Context unexpectedly has no errors after writing JSON with error")
	}
}
