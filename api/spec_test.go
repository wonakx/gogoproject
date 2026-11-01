package apidoc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	apidoc "example.com/gogoproject/api"
	"example.com/gogoproject/internal/api"
	"github.com/getkin/kin-openapi/openapi3"
)

func TestSpecificationMatchesGeneratedContract(t *testing.T) {
	source, err := openapi3.NewLoader().LoadFromData(apidoc.Spec)
	if err != nil {
		t.Fatal(err)
	}
	if err := source.Validate(context.Background()); err != nil {
		t.Fatal(err)
	}
	generated, err := api.GetSpec()
	if err != nil {
		t.Fatal(err)
	}
	want, err := json.Marshal(source)
	if err != nil {
		t.Fatal(err)
	}
	got, err := json.Marshal(generated)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(want, got) {
		t.Fatal("OpenAPI contract changed: run go generate ./api")
	}
}
