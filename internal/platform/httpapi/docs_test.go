package httpapi

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	apidoc "example.com/gogoproject/api"
)

func TestOfflineDocumentation(t *testing.T) {
	h := DocumentationHandler()
	for _, path := range []string{"/docs/", "/openapi.yaml", "/docs/assets/swagger-ui.css", "/docs/assets/swagger-ui-bundle.js"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != 200 || w.Body.Len() == 0 {
			t.Fatalf("%s: %d", path, w.Code)
		}
		if path == "/openapi.yaml" && !bytes.Equal(w.Body.Bytes(), apidoc.Spec) {
			t.Fatal("served contract differs")
		}
		if path == "/docs/" && (!strings.Contains(w.Body.String(), "url:'/openapi.yaml'") || !strings.Contains(w.Body.String(), "supportedSubmitMethods:[]")) {
			t.Fatal("UI not connected to contract")
		}
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("POST", "/api/v1/members", strings.NewReader(`{"id":"x"}`)))
	if w.Code != 404 {
		t.Fatalf("documentation server exposes business API: %d", w.Code)
	}
}
