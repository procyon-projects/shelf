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

	IdField   *FieldMetadata
	ColumnMap map[string]string
	FieldMap  map[string]string
	Fields    []*FieldMetadata
}

func FindEntityTypes(structTypes []marker.StructType) {
	for _, structType := range structTypes {

		var entityMarker shelf.EntityMarker

		if value, ok := StructMustBeMarkedOnceAtMost(structType, shelf.MarkerEntity); ok {
			entityMarker = value.(shelf.EntityMarker)
		} else {
			continue
		}

		metadata := &EntityMetadata{
			EntityName: strings.TrimSpace(structType.Name),
			TableName:  shelf.ToSnakeCase(strings.TrimSpace(structType.Name)),
			StructName: structType.Name,
			StructType: structType,
		}

		if strings.TrimSpace(entityMarker.Name) != "" {
			metadata.EntityName = strings.TrimSpace(entityMarker.Name)
		}

		if _, ok := entitiesByName[metadata.EntityName]; ok {
			err := fmt.Errorf("there is already an entity with name '%s'", metadata.EntityName)
			errs = append(errs, marker.NewError(err, structType.File.FullPath, structType.Position))
			break
		}

		if value, ok := StructMustBeMarkedOnceAtMost(structType, shelf.MarkerTable); ok {
			tableMarker := value.(shelf.TableMarker)

			if strings.TrimSpace(tableMarker.Name) != "" {
				metadata.TableName = strings.TrimSpace(tableMarker.Name)
			}

			if _, ok := entitiesByTableName[metadata.TableName]; ok {
				err := fmt.Errorf("there is already an entity with table name '%s'", metadata.TableName)
				errs = append(errs, marker.NewError(err, structType.File.FullPath, structType.Position))
				break
			}
		}

		fullStructName := FullStructName(structType)
		entityMetadataByStructName[fullStructName] = metadata
		entitiesByTableName[metadata.TableName] = fullStructName
		entitiesByName[metadata.EntityName] = fullStructName
	}
}

func ProcessEntityTypes() {
	for _, metadata := range entityMetadataByStructName {
		ProcessEntityType(metadata)
	}
}

func ProcessEntityType(metadata *EntityMetadata) {
	columns := make(map[string]bool)
	fieldMetadataArr := CollectFieldMetadata(metadata.StructType)

	for _, fieldMetadata := range fieldMetadataArr {

		if fieldMetadata.MarkerFlags&Column == Column {
			if fieldMetadata.ColumnMarker.Name != "" {
				fieldMetadata.ColumnName = fieldMetadata.ColumnMarker.Name
			}
			fieldMetadata.ColumnLength = fieldMetadata.ColumnMarker.Length
			fieldMetadata.UniqueColumn = fieldMetadata.ColumnMarker.Unique
		}

		if _, ok := columns[fieldMetadata.ColumnName]; ok {
			err := fmt.Errorf("there is already a column with name '%s'", fieldMetadata.ColumnName)
			errs = append(errs, marker.NewError(err, fieldMetadata.File.FullPath, fieldMetadata.Position))
			continue
		}

		columns[fieldMetadata.ColumnName] = true
		metadata.ColumnMap[fieldMetadata.ColumnName] = fieldMetadata.FieldName
		metadata.FieldMap[fieldMetadata.FieldName] = fieldMetadata.ColumnName
	}
}
