package openapi3

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolvePath(t *testing.T) {
	var b = &url.URL{Path: "testdata/recursiveRef"}
	var c = &url.URL{Path: "./components/models/error.yaml"}
	u, err := resolvePath(b, c)

	require.NoError(t, err)
	require.Equal(t, &url.URL{Path: "testdata/recursiveRef/components/models/error.yaml"}, u)

	b = &url.URL{Path: "testdata/recursiveRef/openapi.yaml"}
	c = &url.URL{Path: "./components/models/error.yaml"}
	u, err = resolvePath(b, c)

	require.NoError(t, err)
	require.Equal(t, &url.URL{Path: "testdata/recursiveRef/components/models/error.yaml"}, u)
}

func TestLoaderSupportsRecursiveReference(t *testing.T) {
	loader := NewLoader()
	loader.IsExternalRefsAllowed = true
	doc, err := loader.LoadFromFile("testdata/recursiveRef/openapi.yml")
	require.NoError(t, err)
	err = doc.Validate(loader.Context)
	require.NoError(t, err)
	require.Equal(t, "bar", doc.Paths["/foo"].Get.Responses.Get(200).Value.Content.Get("application/json").Schema.Value.Properties["foo2"].Value.Properties["foo"].Value.Properties["bar"].Value.Example)
	require.Equal(t, "ErrorDetails", doc.Paths["/foo"].Get.Responses.Get(400).Value.Content.Get("application/json").Schema.Value.Title)
	require.Equal(t, "ErrorDetails", doc.Paths["/double-ref-foo"].Get.Responses.Get(400).Value.Content.Get("application/json").Schema.Value.Title)
}

func TestIssue447(t *testing.T) {
	loader := NewLoader()
	doc, err := loader.LoadFromData([]byte(`
openapi: 3.0.1
info:
  title: Recursive refs example
  version: "1.0"
paths: {}
components:
  schemas:
    Complex:
      type: object
      properties:
        parent:
          $ref: '#/components/schemas/Complex'
`))
	require.NoError(t, err)
	err = doc.Validate(loader.Context)
	require.NoError(t, err)
	require.Equal(t, "object", doc.Components.
		// Complex
		Schemas["Complex"].
		// parent
		Value.Properties["parent"].
		// parent
		Value.Properties["parent"].
		// parent
		Value.Properties["parent"].
		// type
		Value.Type)
}
