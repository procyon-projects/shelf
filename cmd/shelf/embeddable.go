package main

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/shelf"
)

type EmbeddableMetadata struct {
	StructName string
	StructType marker.StructType

	ColumnMap     map[string]string
	FieldMap      map[string]string
	Fields        []*FieldMetadata
	InheritFields []*FieldMetadata
}

func (metadata *EmbeddableMetadata) FindFieldMetadataByFieldName(fieldName string) (*FieldMetadata, bool) {
	for _, field := range metadata.Fields {
		if field.FieldName == fieldName {
			return field, true
		}
	}

	return nil, false
}

func FindEmbeddableType(structTypes []marker.StructType) {
	for _, structType := range structTypes {
		_, markedAsEmbeddable := MarkedAs(structType.Markers, shelf.MarkerEmbeddable)

		if !markedAsEmbeddable {
			continue
		}

		_, markedAsEntity := MarkedAs(structType.Markers, shelf.MarkerEntity)

		if markedAsEntity {
			err := fmt.Errorf("the struct cannot be marked as both shelf:entity and shelf:embeddable markers")
			errs = append(errs, marker.NewError(err, structType.File.FullPath, structType.Position))
			continue
		}

		embeddableMetadata := &EmbeddableMetadata{
			StructName:    structType.Name,
			StructType:    structType,
			ColumnMap:     make(map[string]string),
			FieldMap:      make(map[string]string),
			Fields:        make([]*FieldMetadata, 0),
			InheritFields: make([]*FieldMetadata, 0),
		}

		embeddableMetadataByStructName[FullStructName(structType)] = embeddableMetadata
	}
}

func ProcessEmbeddableTypes() {
	for _, metadata := range embeddableMetadataByStructName {
		PreProcessEmbeddableType(metadata)
	}

	PostProcessEmbeddableTypes()
}

func PreProcessEmbeddableType(metadata *EmbeddableMetadata) {
	fieldMetadataArr := CollectFieldMetadata(metadata.StructType)

	for _, fieldMetadata := range fieldMetadataArr {
		if fieldMetadata.TypeMetadata.HasAnError {
			continue
		}

		if _, ok := metadata.ColumnMap[fieldMetadata.ColumnName]; ok {
			err := fmt.Errorf("there is already a column with name '%s'", fieldMetadata.ColumnName)
			errs = append(errs, marker.NewError(err, fieldMetadata.File.FullPath, fieldMetadata.Position))
		}

		if fieldMetadata.MarkerFlags&Associations != 0 {
			err := fmt.Errorf("the embeddable types do not support associations fields")
			errs = append(errs, marker.NewError(err, fieldMetadata.File.FullPath, fieldMetadata.Position))
		}

		if !fieldMetadata.IsEmbedded && fieldMetadata.MarkerFlags&Transient == 0 && fieldMetadata.MarkerFlags&Embedded == 0 &&
			fieldMetadata.MarkerFlags&Associations == 0 {
			metadata.ColumnMap[fieldMetadata.ColumnName] = fieldMetadata.FieldName
			metadata.FieldMap[fieldMetadata.FieldName] = fieldMetadata.ColumnName
		}

		metadata.Fields = append(metadata.Fields, fieldMetadata)
	}
}

func PostProcessEmbeddableTypes() {
	for _, metadata := range embeddableMetadataByStructName {
		PostProcessEmbeddableType(metadata)
	}
}

func PostProcessEmbeddableType(metadata *EmbeddableMetadata) {
	for _, fieldMetadata := range metadata.Fields {
		if fieldMetadata.IsEmbedded {
			embeddedMetadata, embeddableExists := embeddableMetadataByStructName[fieldMetadata.TypeMetadata.ImportPath+"#"+fieldMetadata.TypeMetadata.SimpleTypeName]

			if !embeddableExists {
				err := fmt.Errorf("the field type '%s' cannot be embedded because the type is not marked as shelf:embeddable marker", fieldMetadata.TypeMetadata.TypeName)
				errs = append(errs, marker.NewError(err, fieldMetadata.File.FullPath, fieldMetadata.Position))
			} else {

				for _, field := range embeddedMetadata.Fields {
					if inheritFieldMetadata, ok := embeddedMetadata.FindFieldMetadataByFieldName(field.FieldName); ok {
						metadata.InheritFields = append(metadata.InheritFields, inheritFieldMetadata)
					}
				}

			}
		}
	}
}

func FindColumnMap(metadata *EmbeddableMetadata) {
	for _, fieldMetadata := range metadata.Fields {
		if fieldMetadata.IsEmbedded {
			embeddedMetadata, embeddableExists := embeddableMetadataByStructName[fieldMetadata.TypeMetadata.ImportPath+"#"+fieldMetadata.TypeMetadata.SimpleTypeName]

			if !embeddableExists {
				continue
			}

			FindColumnMap(embeddedMetadata)
		} else {
			metadata.InheritFields = append(metadata.InheritFields)
		}
	}

}
