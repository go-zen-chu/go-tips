package httpserver

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test sample using httptest
func Test_handler_GetReq(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/test", bytes.NewBufferString(`{
			"text": "testtest"
			}`))
		w := httptest.NewRecorder()
		hdl := NewHttpHandler()
		hdl.GetReq(w, r)

		if w.Code != http.StatusOK {
			t.Error("should be status OK")
		}
	})
}
