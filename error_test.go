package bowtie

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	m := "Hello there."
	err := errors.New(m)

	// Create a new error from a traditional Go error

	e := NewErrorWithError(err).CaptureStackTrace()

	if e.StatusCode() != 500 {
		t.Errorf("Expected status code 500, got %d instead", e.StatusCode())
	}

	if e.Error() == m {
		t.Errorf("Expected a generic error message, got %s instead", e.Error())
	}

	if e.String() == m {
		t.Errorf("Expected a generic error message, got %s instead", e.Error())
	}

	data, err := json.Marshal(e)

	if err != nil {
		t.Fatalf("Unable to marshal Error instance to JSON: %s", err)
	}

	if string(data) != `{"message":"An server error has occurred.","statusCode":500}` {
		t.Errorf("Unexpected JSON marshal received: %s", string(data))
	}

	if len(e.StackTrace()) != 3 {
		t.Errorf("Unexpected stack trace: %#v", e.StackTrace())
	}
}
