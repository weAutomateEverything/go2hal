package gokit

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestEncodeRequest(t *testing.T) {
	r := http.Request{}
	EncodeRequest(nil, &r, "testing 123")

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Esception in Encode Request %v", err)
	}

	if string(result) != "testing 123" {
		t.Errorf("Unexpected result %v, expected testing 123", string(result))
	}

}
