package main

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/shelf"
)

type FieldMetadata struct {
	FieldName    string
	ColumnName   string
	ColumnLength int
	UniqueColumn bool

	File      *marker.File
	Position  marker.Position
	Field     marker.Field
	FieldType marker.Type

	MarkerFlags      MarkerFlag
	ColumnMarker     shelf.ColumnMarker
	TemporalMarker   shelf.TemporalMarker
	EnumeratedMarker shelf.EnumeratedMarker

	AttributeOverrideMarkers []shelf.AttributeOverrideMarker

	OneToOneMarker   shelf.OneToOneMarker
	OneToManyMarker  shelf.OneToManyMarker
	ManyToOneMarker  shelf.ManyToOneMarker
	ManyToManyMarker shelf.ManyToManyMarker
}

func FieldMetadataFromField(field marker.Field) *FieldMetadata {
	return &FieldMetadata{
		FieldName:  field.Name,
		ColumnName: shelf.ToSnakeCase(field.Name),
		File:       field.File,
		Position:   field.Position,
		Field:      field,
		FieldType:  field.Type,
	}
}

func CollectFieldMetadata(structType marker.StructType) []*FieldMetadata {
	fields := make([]*FieldMetadata, 0)

	for _, field := range structType.Fields {
		metadata := FieldMetadataFromField(field)
		fields = append(fields, metadata)

		if !field.IsExported {
			err := fmt.Errorf("the field '%s' in '%s' must be exported", field.Name, structType.Name)
			errs = append(errs, marker.NewError(err, field.File.FullPath, field.Position))
		}

		if value, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerColumn); ok {
			metadata.MarkerFlags |= Column
			columnMarker := value.(shelf.ColumnMarker)
			field.Name = columnMarker.Name
		}

		if _, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerTransient); ok {
			metadata.MarkerFlags |= Transient
		}

		if _, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerId); ok {
			metadata.MarkerFlags |= Id
		}

		if _, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerGeneratedValue); ok {
			metadata.MarkerFlags |= GeneratedValue
		}

		if _, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerLob); ok {
			metadata.MarkerFlags |= Lob
		}

		if value, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerTemporal); ok {
			metadata.MarkerFlags |= Temporal
			metadata.TemporalMarker = value.(shelf.TemporalMarker)
		}

		if _, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerCreatedDate); ok {
			metadata.MarkerFlags |= CreatedDate
		}

		if _, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerLastModifiedDate); ok {
			metadata.MarkerFlags |= LastModifiedDate
		}

		if _, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerEmbedded); ok {
			metadata.MarkerFlags |= Embedded
		}

		if values, ok := FieldCanBeMarkedManyTimes(field, shelf.MarkerAttributeOverride); ok {
			metadata.MarkerFlags |= AttributeOverride
			metadata.AttributeOverrideMarkers = ToAttributeOverrideMarkers(values)
		}

		if value, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerEnumerated); ok {
			metadata.MarkerFlags |= Enumerated
			metadata.EnumeratedMarker = value.(shelf.EnumeratedMarker)
		}

		if value, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerOneToOne); ok {
			metadata.MarkerFlags |= OneToOne
			metadata.OneToOneMarker = value.(shelf.OneToOneMarker)
		}

		if value, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerOneToMany); ok {
			metadata.MarkerFlags |= OneToMany
			metadata.OneToManyMarker = value.(shelf.OneToManyMarker)
		}

		if value, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerManyToOne); ok {
			metadata.MarkerFlags |= ManyToOne
			metadata.ManyToOneMarker = value.(shelf.ManyToOneMarker)
		}

		if value, ok := FieldMustBeMarkedOnceAtMost(field, shelf.MarkerManyToMany); ok {
			metadata.MarkerFlags |= ManyToMany
			metadata.ManyToManyMarker = value.(shelf.ManyToManyMarker)
		}
	}

	return fields
}
