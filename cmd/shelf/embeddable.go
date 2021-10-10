package main

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/shelf"
)

type EmbeddableMetadata struct {
	StructName string
	StructType marker.StructType
}

func ValidateEmbeddableMarkers(structType marker.StructType) bool {
	markers := structType.Markers

	if markers == nil {
		return true
	}

	isValid := true

	for name, values := range markers {
		if values != nil && len(values) > 1 {
			err := fmt.Errorf("the struct cannot be marked twice as '%s' marker", name)
			errs = append(errs, marker.NewError(err, structType.File.FullPath, marker.Position{
				Line:   structType.Position.Line,
				Column: structType.Position.Column,
			}))
			isValid = false
		}
	}

	return isValid
}

func FindEmbeddables(structTypes []marker.StructType) {
	for _, structType := range structTypes {
		if !ValidateEmbeddableMarkers(structType) {
			continue
		}

		markerValues := structType.Markers

		if markerValues == nil {
			continue
		}

		markers, ok := markerValues[shelf.MarkerEmbeddable]

		if !ok {
			continue
		}

		matched := false

		for _, candidateMarker := range markers {
			switch candidateMarker.(type) {
			case shelf.EmbeddableMarker:
				matched = true
			}
		}

		if matched {
			embeddableMetadata := &EmbeddableMetadata{
				StructName: structType.Name,
				StructType: structType,
			}

			fullStructName := structType.File.Package.Path + "#" + structType.Name
			embeddableMetadataByStructName[fullStructName] = embeddableMetadata
		}
	}
}
