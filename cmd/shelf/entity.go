package main

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/shelf"
	"strings"
)

type EntityMetadata struct {
	EntityName string
	TableName  string
	StructName string
	StructType marker.StructType
}

func ValidateEntityMarkers(structType marker.StructType) bool {
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

func FindEntities(structTypes []marker.StructType) {
	for _, structType := range structTypes {
		if !ValidateEntityMarkers(structType) {
			continue
		}

		markerValues := structType.Markers

		if markerValues == nil {
			continue
		}

		markers, ok := markerValues[shelf.MarkerEntity]

		if !ok {
			return
		}

		var err error
		entityName := strings.TrimSpace(structType.Name)
		tableName := shelf.ToSnakeCase(strings.TrimSpace(structType.Name))

		for _, candidateMarker := range markers {
			switch typedMarker := candidateMarker.(type) {
			case shelf.EntityMarker:
				entityNameValue := strings.TrimSpace(typedMarker.Name)

				if entityNameValue != "" {
					entityName = entityNameValue
				}

				if _, ok := entitiesByName[entityName]; ok {
					err = fmt.Errorf("there is already an entity with name '%s'", entityName)
					errs = append(errs, marker.NewError(err, structType.File.FullPath, marker.Position{
						Line:   structType.Position.Line,
						Column: structType.Position.Column,
					}))
					break
				}

			case shelf.TableMarker:
				tableNameValue := strings.TrimSpace(typedMarker.Name)

				if tableNameValue != "" {
					tableName = tableNameValue
				}

				if _, ok := entitiesByTableName[tableNameValue]; ok {
					err = fmt.Errorf("there is already an entity with table name '%s'", tableNameValue)
					errs = append(errs, marker.NewError(err, structType.File.FullPath, marker.Position{
						Line:   structType.Position.Line,
						Column: structType.Position.Column,
					}))
					break
				}
			}
		}

		if err == nil {
			entityMetadata := EntityMetadata{
				EntityName: entityName,
				TableName:  tableName,
				StructName: structType.Name,
				StructType: structType,
			}

			fullStructName := structType.File.Package.Path + "#" + structType.Name

			entityMetadataByStructName[fullStructName] = entityMetadata
			entitiesByTableName[tableName] = fullStructName
			entitiesByName[entityName] = fullStructName
		}
	}
}
