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

	IdField *FieldMetadata

	AllColumnMap map[string]string
	AllFields    []*FieldMetadata
	FieldKeyMap  map[string]*FieldMetadata

	ColumnMap map[string]string
	FieldMap  map[string]string
	Fields    []*FieldMetadata
}

func (metadata *EntityMetadata) FindFieldMetadataByFieldName(fieldName string) (*FieldMetadata, bool) {
	for _, field := range metadata.Fields {
		if field.FieldName == fieldName {
			return field, true
		}
	}

	return nil, false
}

func FindEntityTypes(structTypes []marker.StructType) {
	for _, structType := range structTypes {

		var entityMarker shelf.EntityMarker

		if value, ok := StructMustBeMarkedAtMostOnce(structType, shelf.MarkerEntity); ok {
			entityMarker = value.(shelf.EntityMarker)
		} else {
			continue
		}

		metadata := &EntityMetadata{
			EntityName:   strings.TrimSpace(structType.Name),
			TableName:    shelf.ToSnakeCase(strings.TrimSpace(structType.Name)),
			StructName:   structType.Name,
			StructType:   structType,
			AllColumnMap: make(map[string]string),
			ColumnMap:    make(map[string]string),
			FieldMap:     make(map[string]string),
			FieldKeyMap:  make(map[string]*FieldMetadata),
			Fields:       make([]*FieldMetadata, 0),
		}

		if strings.TrimSpace(entityMarker.Name) != "" {
			metadata.EntityName = strings.TrimSpace(entityMarker.Name)
		}

		if _, ok := entitiesByName[metadata.EntityName]; ok {
			err := fmt.Errorf("there is already an entity with name '%s'", metadata.EntityName)
			errs = append(errs, marker.NewError(err, structType.File.FullPath, structType.Position))
			break
		}

		if value, ok := StructMustBeMarkedAtMostOnce(structType, shelf.MarkerTable); ok {
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
		PreProcessEntityType(metadata)
	}

	PostProcessEntityTypes()
}

func PreProcessEntityType(metadata *EntityMetadata) {
	fieldMetadataArr := CollectFieldMetadata(metadata.StructType)

	for _, fieldMetadata := range fieldMetadataArr {
		if fieldMetadata.TypeMetadata.HasAnError {
			continue
		}

		if _, ok := metadata.ColumnMap[fieldMetadata.ColumnName]; ok {
			err := fmt.Errorf("there is already a column with name '%s'", fieldMetadata.ColumnName)
			errs = append(errs, marker.NewError(err, fieldMetadata.File.FullPath, fieldMetadata.Position))
		}

		if !fieldMetadata.IsEmbedded && fieldMetadata.MarkerFlags&Transient == 0 && fieldMetadata.MarkerFlags&Embedded == 0 &&
			fieldMetadata.MarkerFlags&Associations == 0 {
			metadata.ColumnMap[fieldMetadata.ColumnName] = fieldMetadata.FieldName
			metadata.FieldMap[fieldMetadata.FieldName] = fieldMetadata.ColumnName
		}

		if fieldMetadata.IsEmbedded {
			_, embeddableExists := embeddableMetadataByStructName[fieldMetadata.TypeMetadata.ImportPath+"#"+fieldMetadata.TypeMetadata.SimpleTypeName]

			if !embeddableExists {
				err := fmt.Errorf("the field type '%s' cannot be embedded because the type is not marked as shelf:embeddable marker", fieldMetadata.TypeMetadata.TypeName)
				errs = append(errs, marker.NewError(err, fieldMetadata.File.FullPath, fieldMetadata.Position))
			}
		}

		metadata.Fields = append(metadata.Fields, fieldMetadata)
	}
}

func PostProcessEntityTypes() {
	for _, metadata := range entityMetadataByStructName {
		PostProcessEntityType(metadata)
	}
}

func PostProcessEntityType(metadata *EntityMetadata) {
	fieldKeyMap := make(map[string]*FieldMetadata, 0)
	metadata.AllFields = CollectAllFields("", metadata.Fields, nil, fieldKeyMap, nil)
	metadata.FieldKeyMap = fieldKeyMap

	for _, field := range metadata.AllFields {

		if !field.IsEmbedded && field.MarkerFlags&Transient == 0 && field.MarkerFlags&Embedded == 0 &&
			field.MarkerFlags&Associations == 0 {

			if _, ok := metadata.AllColumnMap[field.ColumnName]; ok {
				err := fmt.Errorf("there is already a column with name '%s'", field.ColumnName)
				position := field.Position

				if field.ParentField != nil {
					position = field.ParentField.Position
				}

				errs = append(errs, marker.NewError(err, field.File.FullPath, position))
			} else {
				metadata.AllColumnMap[field.ColumnName] = field.FieldName
			}

		}
	}
}
