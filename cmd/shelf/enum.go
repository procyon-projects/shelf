package main

import (
	"github.com/procyon-projects/marker"
	"strings"
)

type EnumMetadata struct {
	Name            string
	UserDefinedType marker.UserDefinedType
	Values          []marker.ConstValue
}

func FindEnums(file *marker.File) {
	for _, typ := range file.UserDefinedTypes {

		enumMetadata := &EnumMetadata{
			Name:            typ.Name,
			UserDefinedType: typ,
		}

		canBeEnum := false

		switch typedActualType := typ.ActualType.(type) {
		case *marker.ObjectType:
			if typedActualType.ImportName == "" && (strings.HasPrefix(typedActualType.Name, "int") ||
				strings.HasPrefix(typedActualType.Name, "uint")) {
				canBeEnum = true
			}
		}

		if !canBeEnum {
			continue
		}

		fullEnumName := typ.File.Package.Path + "#" + typ.Name

		for _, constValue := range file.Consts {
			valueType := constValue.Type

			if valueType != nil && valueType.ImportName == "" && valueType.Name == typ.Name {
				enumMetadata.Values = append(enumMetadata.Values, constValue)
			}
		}

		enumsByName[fullEnumName] = enumMetadata
	}

}
