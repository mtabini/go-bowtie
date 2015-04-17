package bowtie

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Struct Request adds a few convenience functions to `http.Request`.
type Request struct {
	*http.Request
}

// NewRequest creates a new request instance. This is called transparently for you
// at the time the server receives a request
func NewRequest(r *http.Request) *Request {
	return &Request{r}
}

// StringBody returns the request's body as a string
func (r *Request) StringBody() (string, error) {
	if r.Body != nil {
		res, err := ioutil.ReadAll(r.Body)

		return string(res), err
	}

	return "", nil
}

// JSONBody attempts to unmarshal JSON out of the request's body, and
// returns a map if successful, or an error if not.
func (r *Request) JSONBody() (map[string]interface{}, error) {
	if r.Body != nil {
		res := map[string]interface{}{}

		err := json.NewDecoder(r.Body).Decode(&res)

		return res, err
	}

	return map[string]interface{}{}, nil
}

// ReadJSONBody attempts to unmarshal JSON from the request's body into
// a destination of your choosing.
func (r *Request) ReadJSONBody(v interface{}) error {
	if r.Body != nil {
		err := json.NewDecoder(r.Body).Decode(&v)

		return err
	}

	return nil
}
