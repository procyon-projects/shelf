/*
Copyright © 2021 Shelf Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/procyon-projects/marker"
	"log"
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

// printErrors prints error.
func PrintError(err error) {
	if err != nil {

		switch typedErr := err.(type) {
		case marker.ErrorList:
			PrintErrors(typedErr)
			return
		}

		log.Errorln(err)
		return
	}
}

// printErrors prints the error list.
func PrintErrors(errorList marker.ErrorList) {
	if errorList == nil || len(errorList) == 0 {
		return
	}

	for _, err := range errorList {
		switch typedErr := err.(type) {
		case marker.Error:
			pos := typedErr.Position
			log.Errorf("%s (%d:%d) : %s\n", typedErr.FileName, pos.Line, pos.Column, typedErr.Error())
		case marker.ParserError:
			pos := typedErr.Position
			log.Errorf("%s (%d:%d) : %s\n", typedErr.FileName, pos.Line, pos.Column, typedErr.Error())
		case marker.ErrorList:
			PrintErrors(typedErr)
		default:
			PrintError(err)
		}
	}
}

// validateMarkers visits all files and returns errors
func ValidateMarkers(collector *marker.Collector, pkgs []*marker.Package) error {
	marker.EachFile(collector, pkgs, func(file *marker.File, fileErrors error) {
		if fileErrors != nil {
			validationErrors = append(validationErrors, fileErrors)
		}
	})

	return marker.NewErrorList(validationErrors)
}

func TypeNameFromTypeWithoutImport(typ marker.Type) string {
	switch typed := typ.(type) {
	case *marker.ObjectType:
		return typed.Name
	case *marker.PointerType:
		return TypeNameFromType(typed.Typ)
	}

	return ""
}

func TypeNameFromType(typ marker.Type) string {
	switch typed := typ.(type) {
	case *marker.ObjectType:
		name := typed.Name

		if typed.ImportName != "" {
			name = typed.ImportName + "." + name
		}
		return name
	case *marker.PointerType:
		return TypeNameFromType(typed.Typ)
	}

	return ""
}

func GetObjectType(typ marker.Type) *marker.ObjectType {
	switch typed := typ.(type) {
	case *marker.ObjectType:
		return typed
	case *marker.PointerType:
		return GetObjectType(typed.Typ)
	case *marker.ArrayType:
		return GetObjectType(typed.ItemType)
	default:
		return nil
	}
}

func IsArrayType(typ marker.Type) bool {
	switch typed := typ.(type) {
	case *marker.ArrayType:
		return true
	case *marker.PointerType:
		return IsArrayType(typed.Typ)
	}

	return false
}

func IsBasicType(typ marker.Type) (ok, supported bool) {
	var objectType *marker.ObjectType

	switch typed := typ.(type) {
	case *marker.ObjectType:
		objectType = typed
	case *marker.PointerType:
		return IsBasicType(typed.Typ)
	default:
		return
	}

	if objectType.ImportName != "" {
		return
	}

	if isSupported, exists := basicTypes[objectType.Name]; exists {
		supported = isSupported
		ok = true
	}

	return
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

		if typ.ImportName == fileImport.Path {
			name = fileImport.Path
			path = fileImport.Path
			exists = true
			return
		}

	}

	return
}
