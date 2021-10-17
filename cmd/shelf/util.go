package main

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/shelf"
	"strings"
)

const NameSeparator = "#"

func FullStructName(structType marker.StructType) string {
	return structType.File.Package.Path + NameSeparator + structType.Name
}

func FullInterfaceName(interfaceType marker.InterfaceType) string {
	return interfaceType.File.Package.Path + NameSeparator + interfaceType.Name
}

func MarkedAs(values marker.MarkerValues, marker string) ([]interface{}, bool) {
	if values == nil {
		return nil, false
	}

	markers, ok := values[marker]

	if !ok {
		return nil, false
	}

	return markers, true
}

func MustBeMarkerAtMostOnce(typeName string, markers marker.MarkerValues, markerName string, file *marker.File, position marker.Position) (interface{}, bool) {
	values, ok := MarkedAs(markers, markerName)

	if !ok {
		return nil, false
	}

	if len(values) > 1 {
		err := fmt.Errorf("the '%s' cannot be marked twice as '%s' marker", typeName, markerName)
		errs = append(errs, marker.NewError(err, file.FullPath, position))
	}

	return values[0], true
}

func StructMustBeMarkedAtMostOnce(structType marker.StructType, markerName string) (interface{}, bool) {
	return MustBeMarkerAtMostOnce(structType.Name, structType.Markers, markerName, structType.File, structType.Position)
}

func InterfaceMustBeMarkedAtMostOnce(interfaceType marker.InterfaceType, markerName string) (interface{}, bool) {
	return MustBeMarkerAtMostOnce(interfaceType.Name, interfaceType.Markers, markerName, interfaceType.File, interfaceType.Position)
}

func FieldMustBeMarkedAtMostOnce(field marker.Field, markerName string) (interface{}, bool) {
	return MustBeMarkerAtMostOnce(field.Name, field.Markers, markerName, field.File, field.Position)
}

func FieldCanBeMarkedManyTimes(field marker.Field, markerName string) ([]interface{}, bool) {
	values, ok := MarkedAs(field.Markers, markerName)

	if !ok {
		return nil, false
	}

	return values, true
}

func CanBeEnum(typ marker.Type) bool {
	switch result := typ.(type) {
	case *marker.ObjectType:
		if result.ImportName == "" && (strings.HasPrefix(result.Name, "int") || strings.HasPrefix(result.Name, "uint")) {
			return true
		}
	}

	return false
}

func ToAttributeOverrideMarkers(markers []interface{}) []shelf.AttributeOverrideMarker {
	attributeOverrideMarkers := make([]shelf.AttributeOverrideMarker, 0)

	for _, element := range markers {
		attributeOverrideMarkers = append(attributeOverrideMarkers, element.(shelf.AttributeOverrideMarker))
	}

	return attributeOverrideMarkers
}
