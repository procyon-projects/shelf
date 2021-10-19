package main

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/shelf"
	"strings"
)

type AttributeOverrideMetadata struct {
	FieldName    string
	ColumnName   string
	ColumnLength int
	UniqueColumn bool
}

type FieldMetadata struct {
	Key         string
	ParentField *FieldMetadata

	FieldName    string
	ColumnName   string
	ColumnLength int
	UniqueColumn bool

	IsEmbedded          bool
	TypeMetadata        *TypeMetadata
	OverrideMetadataMap map[string]*AttributeOverrideMetadata

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
		FieldName:           field.Name,
		ColumnName:          shelf.ToSnakeCase(field.Name),
		IsEmbedded:          field.IsEmbedded,
		TypeMetadata:        TypeMetadataFromType(field.File, field.Type),
		OverrideMetadataMap: make(map[string]*AttributeOverrideMetadata),
		File:                field.File,
		Position:            field.Position,
		Field:               field,
		FieldType:           field.Type,
	}
}

func PreValidateFieldMetadata(metadata *FieldMetadata) {
	if metadata.TypeMetadata.HasAnError {
		var err error
		if metadata.TypeMetadata.TypeName != "" {
			err = fmt.Errorf("the field type '%s' is neither imported nor valid", metadata.TypeMetadata.TypeName)
		} else {
			err = fmt.Errorf("the field type is not supported")
		}
		errs = append(errs, marker.NewError(err, metadata.File.FullPath, metadata.Position))
	}

	if metadata.MarkerFlags&Transient == Transient && metadata.MarkerFlags|Transient != Transient {
		err := fmt.Errorf("the field '%s' with shelf:transient marker cannot have another marker", metadata.FieldName)
		errs = append(errs, marker.NewError(err, metadata.File.FullPath, metadata.Position))
	}

	if metadata.MarkerFlags&Enumerated == Enumerated {
		if metadata.MarkerFlags|Column != (Column | Enumerated) {
			err := fmt.Errorf("the field '%s' with shelf:enumerated marker can only have shelf:column marker", metadata.FieldName)
			errs = append(errs, marker.NewError(err, metadata.File.FullPath, metadata.Position))
		}

		if metadata.TypeMetadata.IsArray {
			err := fmt.Errorf("the enum arrays is not supported")
			errs = append(errs, marker.NewError(err, metadata.File.FullPath, metadata.Position))
		} else {
			_, enumExists := enumsByName[metadata.TypeMetadata.ImportPath+"#"+metadata.TypeMetadata.SimpleTypeName]

			if !enumExists {
				err := fmt.Errorf("the field '%s' cannot be marked as shelf:enumerated because the type is not an enum", metadata.TypeMetadata.TypeName)
				errs = append(errs, marker.NewError(err, metadata.File.FullPath, metadata.Position))
			}
		}
	}

	if metadata.MarkerFlags&Embedded == Embedded {
		if metadata.MarkerFlags|AttributeOverride != (Embedded | AttributeOverride) {
			err := fmt.Errorf("the field '%s' with shelf:embedded marker can only have shelf:attribute-override markers", metadata.FieldName)
			errs = append(errs, marker.NewError(err, metadata.File.FullPath, metadata.Position))
		}

		if metadata.TypeMetadata.IsArray {
			err := fmt.Errorf("the embedded arrays is not supported")
			errs = append(errs, marker.NewError(err, metadata.File.FullPath, metadata.Position))
		} else {
			_, embeddableExists := embeddableMetadataByStructName[metadata.TypeMetadata.ImportPath+"#"+metadata.TypeMetadata.SimpleTypeName]

			if !embeddableExists {
				err := fmt.Errorf("the field '%s' cannot be marked as shelf:embedded because the type is not an embeddedable", metadata.TypeMetadata.TypeName)
				errs = append(errs, marker.NewError(err, metadata.File.FullPath, metadata.Position))
			}
		}
	}
}

func CollectFieldMetadata(structType marker.StructType) []*FieldMetadata {
	fields := make([]*FieldMetadata, 0)

	for _, field := range structType.Fields {
		metadata := FieldMetadataFromField(field)
		fields = append(fields, metadata)

		if !field.IsEmbedded && !field.IsExported {
			err := fmt.Errorf("the field '%s' in '%s' must be exported", field.Name, structType.Name)
			errs = append(errs, marker.NewError(err, field.File.FullPath, field.Position))
		}

		if value, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerColumn); ok {
			metadata.MarkerFlags |= Column
			columnMarker := value.(shelf.ColumnMarker)
			field.Name = columnMarker.Name
		}

		if _, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerTransient); ok {
			metadata.MarkerFlags |= Transient
		}

		if _, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerId); ok {
			metadata.MarkerFlags |= Id
		}

		if _, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerGeneratedValue); ok {
			metadata.MarkerFlags |= GeneratedValue
		}

		if _, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerLob); ok {
			metadata.MarkerFlags |= Lob
		}

		if value, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerTemporal); ok {
			metadata.MarkerFlags |= Temporal
			metadata.TemporalMarker = value.(shelf.TemporalMarker)
		}

		if _, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerCreatedDate); ok {
			metadata.MarkerFlags |= CreatedDate
		}

		if _, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerLastModifiedDate); ok {
			metadata.MarkerFlags |= LastModifiedDate
		}

		if _, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerEmbedded); ok {
			metadata.MarkerFlags |= Embedded
		}

		if values, ok := FieldCanBeMarkedManyTimes(field, shelf.MarkerAttributeOverride); ok {
			metadata.MarkerFlags |= AttributeOverride
			metadata.AttributeOverrideMarkers = ToAttributeOverrideMarkers(values)
		}

		if value, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerEnumerated); ok {
			metadata.MarkerFlags |= Enumerated
			metadata.EnumeratedMarker = value.(shelf.EnumeratedMarker)
		}

		if value, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerOneToOne); ok {
			metadata.MarkerFlags |= OneToOne
			metadata.OneToOneMarker = value.(shelf.OneToOneMarker)
		}

		if value, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerOneToMany); ok {
			metadata.MarkerFlags |= OneToMany
			metadata.OneToManyMarker = value.(shelf.OneToManyMarker)
		}

		if value, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerManyToOne); ok {
			metadata.MarkerFlags |= ManyToOne
			metadata.ManyToOneMarker = value.(shelf.ManyToOneMarker)
		}

		if value, ok := FieldMustBeMarkedAtMostOnce(field, shelf.MarkerManyToMany); ok {
			metadata.MarkerFlags |= ManyToMany
			metadata.ManyToManyMarker = value.(shelf.ManyToManyMarker)
		}

		PreValidateFieldMetadata(metadata)

		if metadata.MarkerFlags&Column == Column {
			if metadata.ColumnMarker.Name != "" {
				metadata.ColumnName = metadata.ColumnMarker.Name
			}
			metadata.ColumnLength = metadata.ColumnMarker.Length
			metadata.UniqueColumn = metadata.ColumnMarker.Unique
		}

		if !metadata.IsEmbedded && metadata.MarkerFlags&Embedded == Embedded && metadata.MarkerFlags&AttributeOverride == AttributeOverride {
			for _, overrideMarker := range metadata.AttributeOverrideMarkers {
				metadata.OverrideMetadataMap[overrideMarker.Name] = &AttributeOverrideMetadata{
					FieldName:    overrideMarker.Name,
					ColumnName:   overrideMarker.ColumnName,
					ColumnLength: overrideMarker.ColumnLength,
					UniqueColumn: overrideMarker.ColumnUnique,
				}
			}
		}
	}

	return fields
}

func CollectAllFields(key string, fields []*FieldMetadata, parentField *FieldMetadata, fieldKeyMap map[string]*FieldMetadata, overrideMap map[string]*AttributeOverrideMetadata) []*FieldMetadata {
	collectedFields := make([]*FieldMetadata, 0)

	for _, field := range fields {

		fieldKey := strings.Join([]string{key, field.FieldName}, ".")

		if key == "" {
			fieldKey = field.FieldName
		}

		field.Key = fieldKey
		field.ParentField = parentField

		if fieldKey != "" {
			fieldKeyMap[fieldKey] = field
		}

		if field.IsEmbedded || field.MarkerFlags&Embedded == Embedded {
			embeddedMetadata, embeddableExists := embeddableMetadataByStructName[field.TypeMetadata.ImportPath+"#"+field.TypeMetadata.SimpleTypeName]

			if embeddableExists {

				if overrideMap != nil {
					mergedOverrideMap := MergeOverrideMetadataMap(overrideMap, field.FieldName, field.OverrideMetadataMap)
					collectedFields = append(collectedFields, CollectAllFields(fieldKey, embeddedMetadata.Fields, field, fieldKeyMap, mergedOverrideMap)...)
				} else {
					collectedFields = append(collectedFields, CollectAllFields(fieldKey, embeddedMetadata.Fields, field, fieldKeyMap, field.OverrideMetadataMap)...)
				}

			}
			continue
		}

		keys := strings.SplitN(key, ".", 2)

		var overrideMetadataKey = ""

		if len(keys) == 1 && keys[0] != "" {
			overrideMetadataKey = field.FieldName
		} else if len(keys) == 2 {
			overrideMetadataKey = strings.Join([]string{keys[1], field.FieldName}, ".")
		} else {
			overrideMetadataKey = field.FieldName
		}

		if overrideMap != nil && overrideMetadataKey != "" {
			if overrideMetadata, ok := overrideMap[overrideMetadataKey]; ok {
				field.ColumnName = overrideMetadata.ColumnName
				field.ColumnLength = overrideMetadata.ColumnLength
				field.UniqueColumn = overrideMetadata.UniqueColumn
			} else {
				field.ColumnName = shelf.ToSnakeCase(fieldKey)
			}
		}

		collectedFields = append(collectedFields, field)
	}

	return collectedFields
}

func MergeOverrideMetadataMap(parentOverrideMap map[string]*AttributeOverrideMetadata,
	fieldPrefix string,
	fieldOverrideMap map[string]*AttributeOverrideMetadata) map[string]*AttributeOverrideMetadata {
	result := make(map[string]*AttributeOverrideMetadata, 0)

	if fieldOverrideMap != nil {
		for key, value := range fieldOverrideMap {
			result[fieldPrefix+"."+key] = value
		}
	}

	if parentOverrideMap != nil {
		for key, value := range parentOverrideMap {
			result[key] = value
		}
	}

	return result
}
