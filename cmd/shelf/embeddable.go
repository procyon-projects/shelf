package main

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/shelf"
)

type EmbeddableMetadata struct {
	StructName string
	StructType marker.StructType

	ColumnMap map[string]string
	FieldMap  map[string]string
	Fields    []*FieldMetadata
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
			StructName: structType.Name,
			StructType: structType,
			ColumnMap:  make(map[string]string),
			FieldMap:   make(map[string]string),
			Fields:     make([]*FieldMetadata, 0),
		}

		embeddableMetadataByStructName[FullStructName(structType)] = embeddableMetadata
	}
}

func ProcessEmbeddableTypes() {
	for _, metadata := range embeddableMetadataByStructName {
		ProcessEmbeddableType(metadata)
	}
}

func ProcessEmbeddableType(metadata *EmbeddableMetadata) {
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
