package main

import (
	"github.com/procyon-projects/marker"
	"strings"
)

var basicTypes = map[string]bool{
	"bool":       true,
	"string":     true,
	"int":        true,
	"int8":       true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"uint":       true,
	"uint8":      true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uintptr":    false,
	"byte":       true,
	"rune":       true,
	"float32":    true,
	"float64":    true,
	"complex64":  false,
	"complex128": false,
}

type TypeMetadata struct {
	Type           marker.Type
	ActualType     marker.Type
	SimpleTypeName string
	TypeName       string
	IsBasicType    bool
	IsPointer      bool
	IsArray        bool
	ImportName     string
	ImportPath     string
	HasAnError     bool
}

func TypeMetadataFromType(file *marker.File, typ marker.Type) *TypeMetadata {
	typeMetadata := &TypeMetadata{
		ActualType: typ,
	}

	if pointerType, ok := IsPointerType(typ); ok {
		typeMetadata.IsPointer = true
		typ = pointerType.Typ
	}

	if arrayType, ok := IsArrayType(typ); ok {
		typeMetadata.IsArray = true
		typ = arrayType.ItemType
	} else if objectType, ok := IsObjectType(typ); ok {
		typ = objectType
	}

	if objectType, ok := IsObjectType(typ); ok {
		typeMetadata.Type = objectType
		typeMetadata.SimpleTypeName = objectType.Name

		if strings.TrimSpace(objectType.ImportName) != "" {
			typeMetadata.TypeName = objectType.ImportName + "." + objectType.Name
		} else {
			typeMetadata.TypeName = objectType.Name
		}

		if ok, supported := IsBasicType(objectType); ok {
			if !supported {
				typeMetadata.HasAnError = true
				return typeMetadata
			}

			typeMetadata.IsBasicType = true
		} else {
			name, path, exists := FindImport(file, objectType)

			if exists {
				typeMetadata.ImportName = name
				typeMetadata.ImportPath = path
			} else {
				typeMetadata.HasAnError = true
				return typeMetadata
			}
		}

	} else {
		typeMetadata.HasAnError = true
		return typeMetadata
	}

	return typeMetadata
}

func FindImport(file *marker.File, typ *marker.ObjectType) (name, path string, exists bool) {
	if typ.ImportName == "" {
		path = file.Package.Path
		exists = true
		return
	}

	for _, fileImport := range file.Imports {

		if typ.ImportName == fileImport.Name {
			name = fileImport.Name
			path = fileImport.Path
			exists = true
			return
		}

		if typ.ImportName == fileImport.Path || strings.HasSuffix(fileImport.Path, "/"+typ.ImportName) {
			name = fileImport.Path
			path = fileImport.Path
			exists = true
			return
		}

	}

	return
}

func IsBasicType(objectType *marker.ObjectType) (ok, supported bool) {
	if objectType.ImportName != "" {
		return
	}

	if isSupported, exists := basicTypes[objectType.Name]; exists {
		supported = isSupported
		ok = true
	}

	return
}

func IsObjectType(typ marker.Type) (*marker.ObjectType, bool) {
	switch typed := typ.(type) {
	case *marker.ObjectType:
		return typed, true
	}

	return nil, false
}

func IsPointerType(typ marker.Type) (*marker.PointerType, bool) {
	switch typed := typ.(type) {
	case *marker.PointerType:
		return typed, true
	}

	return nil, false
}

func IsArrayType(typ marker.Type) (*marker.ArrayType, bool) {
	switch typed := typ.(type) {
	case *marker.ArrayType:
		return typed, true
	}

	return nil, false
}
