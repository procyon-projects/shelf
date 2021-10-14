package main

import (
	"github.com/procyon-projects/marker"
)

type EnumMetadata struct {
	Name            string
	UserDefinedType marker.UserDefinedType
	Values          []marker.ConstValue
}

func FindEnumTypes(file *marker.File) {
	for _, typ := range file.UserDefinedTypes {

		if !CanBeEnum(typ.ActualType) {
			continue
		}

		enumMetadata := &EnumMetadata{
			Name:            typ.Name,
			UserDefinedType: typ,
		}

		fullEnumName := typ.File.Package.Path + NameSeparator + typ.Name

		for _, constValue := range file.Consts {
			valueType := constValue.Type

			if valueType == nil {
				continue
			}

			if valueType.ImportName == "" && valueType.Name == typ.Name {
				enumMetadata.Values = append(enumMetadata.Values, constValue)
			}
		}

		enumsByName[fullEnumName] = enumMetadata
	}
}
