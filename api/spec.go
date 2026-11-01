// Package apidoc contains the human-maintained API contract.
package apidoc

import _ "embed"

//go:generate go tool oapi-codegen -config oapi-codegen.yaml -o ../internal/api/api.gen.go openapi.yaml

//go:embed openapi.yaml
var Spec []byte
