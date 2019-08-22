package main

import (
	"github.com/CtrlZvi/jsonschema-minimal-repro/internal/filesystem"
	jsonschema "github.com/xeipuuv/gojsonschema"
)

func main() {
	configurationSchema := jsonschema.NewReferenceLoaderFileSystem(
		"file://schemas/configuration_schema.json",
		filesystem.HTTP,
	)
	_, err := jsonschema.NewSchema(configurationSchema)
	if err != nil {
		panic(err)
	}
}
