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

	IdField *EntityField
	Fields  []*EntityField
}

type EntityField struct {
	FieldName    string
	ColumnName   string
	ColumnLength int
	UniqueColumn bool
	IsTransient  bool
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

func ValidateEntityField(field marker.Field) bool {
	markers := field.Markers

	if markers == nil {
		return true
	}

	isValid := true

	for name, values := range markers {
		if values != nil && len(values) > 1 && name != shelf.MarkerAttributeOverride {
			err := fmt.Errorf("the field cannot be marked twice as '%s' marker", name)
			errs = append(errs, marker.NewError(err, field.File.FullPath, marker.Position{
				Line:   field.Position.Line,
				Column: field.Position.Column,
			}))
			isValid = false
		}
	}

	return isValid
}

func ValidateEntityFields(fields []marker.Field) bool {
	isValid := true
	for _, field := range fields {
		if !ValidateEntityField(field) {
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

		if !ValidateEntityFields(structType.Fields) {
			continue
		}

		markerValues := structType.Markers

		if markerValues == nil {
			continue
		}

		markers, ok := markerValues[shelf.MarkerEntity]

		if !ok {
			continue
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
			entityMetadata := &EntityMetadata{
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

func ProcessEntities() {
	for _, metadata := range entityMetadataByStructName {
		collectFieldMetadata(metadata)
	}
}

func collectFieldMetadata(entityMetadata *EntityMetadata) {
	columns := make(map[string]bool, 0)

	for _, field := range entityMetadata.StructType.Fields {
		position := marker.Position{
			Line:   field.Position.Line,
			Column: field.Position.Column,
		}

		if !field.IsExported {
			err := fmt.Errorf("the field '%s' in '%s' must be exported", field.Name, entityMetadata.StructType.Name)
			errs = append(errs, marker.NewError(err, field.File.FullPath, position))
		}

		if !ValidateFieldType(field, field.Type) {
			continue
		}

		isIdField := false
		hasAnyMarker := false
		hasAnyDateMarker := false
		hasAnyAssociationMarker := false
		hasAssociationError := false

		fieldType := GetObjectType(field.Type)
		isBasicType, isSupported := IsBasicType(field.Type)
		isArrayType := IsArrayType(field.Type)
		_, importPath, ok := FindImport(field.File, fieldType)

		if isBasicType && !isSupported {
			err := fmt.Errorf("the type '%s' is not supported for entities", TypeNameFromType(field.Type))
			errs = append(errs, marker.NewError(err, field.File.FullPath, position))
			continue
		}

		entityField := &EntityField{
			FieldName:  field.Name,
			ColumnName: shelf.ToSnakeCase(field.Name),
		}

		_, hasTransientMarker := field.Markers[shelf.MarkerTransient]

		if hasTransientMarker {
			entityField.IsTransient = true
		}

		_, hasIdMarker := field.Markers[shelf.MarkerId]

		if hasIdMarker {
			if entityMetadata.IdField != nil {
				err := fmt.Errorf("the entity '%s' cannot have more than one Id", entityMetadata.EntityName)
				errs = append(errs, marker.NewError(err, field.File.FullPath, position))
			}

			isIdField = true
			hasAnyMarker = true
		}

		_, hasGeneratedValueMarker := field.Markers[shelf.MarkerGeneratedValue]

		if hasGeneratedValueMarker {
			hasAnyMarker = true
			if !hasIdMarker {
				err := fmt.Errorf("id field can be marked as shelf:generated-value")
				errs = append(errs, marker.NewError(err, field.File.FullPath, position))
			}
		}

		_, hasTemporalMarker := field.Markers[shelf.MarkerTemporal]

		if hasTemporalMarker {
			hasAnyMarker = true
			hasAnyDateMarker = true
		}

		_, hasCreatedDateMarker := field.Markers[shelf.MarkerCreatedDate]

		if hasCreatedDateMarker {
			hasAnyMarker = true
			hasAnyDateMarker = true
		}

		_, hasModifiedDateMarker := field.Markers[shelf.MarkerLastModifiedDate]

		if hasModifiedDateMarker {
			hasAnyMarker = true
			hasAnyDateMarker = true
		}

		_, hasOneToOneMarker := field.Markers[shelf.MarkerOneToOne]

		if hasOneToOneMarker {
			hasAnyAssociationMarker = true
		}

		_, hasOneToManyMarker := field.Markers[shelf.MarkerOneToMany]

		if hasOneToManyMarker {
			if hasAnyAssociationMarker {
				hasAssociationError = true
			}
			hasAnyAssociationMarker = true
		}

		_, hasManyToOne := field.Markers[shelf.MarkerManyToOne]

		if hasManyToOne {
			if hasAnyAssociationMarker {
				hasAssociationError = true
			}
			hasAnyAssociationMarker = true
		}

		_, hasManyToMany := field.Markers[shelf.MarkerManyToMany]

		if hasManyToMany {
			if hasAnyAssociationMarker {
				hasAssociationError = true
			}
			hasAnyAssociationMarker = true
		}

		if hasAssociationError {
			err := fmt.Errorf("the field '%s' must only have at most one of the following association markers: \n"+
				"shelf:one-to-one, shelf:one-to-many, shelf:many-to-many, shelf:many-to-one", field.Name)
			errs = append(errs, marker.NewError(err, field.File.FullPath, position))
		} else if isArrayType && (hasOneToOneMarker || hasManyToOne) {
			err := fmt.Errorf("the array types can only be marked as either shelf:one-to-many or shelf:many-to-many marker")
			errs = append(errs, marker.NewError(err, field.File.FullPath, position))
		} else if !isArrayType && (hasOneToManyMarker || hasManyToMany) {
			err := fmt.Errorf("shelf:one-to-many and shelf:many-to-many markers can only be used for the array types")
			errs = append(errs, marker.NewError(err, field.File.FullPath, position))
		}

		if hasAnyAssociationMarker {
			if hasAnyMarker {
				err := fmt.Errorf("the field '%s' with association marker can only have shelf:column marker", field.Name)
				errs = append(errs, marker.NewError(err, field.File.FullPath, position))
			}
			hasAnyMarker = true
		}

		_, hasEnumeratedMarker := field.Markers[shelf.MarkerEnumerated]

		if hasEnumeratedMarker {
			if hasAnyMarker {
				err := fmt.Errorf("the field '%s' with shelf:enumerated marker can only have shelf:column marker", field.Name)
				errs = append(errs, marker.NewError(err, field.File.FullPath, position))
			}

			hasAnyMarker = true

			if !ok {
				err := fmt.Errorf("the field type '%s' is not imported", TypeNameFromType(field.Type))
				errs = append(errs, marker.NewError(err, field.File.FullPath, position))
			}

			enumName := importPath + "#" + TypeNameFromType(field.Type)
			_, exists := enumsByName[enumName]

			if !exists {
				err := fmt.Errorf("the field '%s' cannot be marked as shelf:enumerated becuase the field type is not an enum", field.Name)
				errs = append(errs, marker.NewError(err, field.File.FullPath, position))
			} else if isArrayType {
				err := fmt.Errorf("the enum arrays is not supported")
				errs = append(errs, marker.NewError(err, field.File.FullPath, position))
			}
		}

		columnMarkers, hasColumnMarker := field.Markers[shelf.MarkerColumn]

		if hasColumnMarker {
			hasAnyMarker = true
			columnMarker := columnMarkers[0].(shelf.ColumnMarker)
			if columnMarker.Name != "" {
				entityField.ColumnName = columnMarker.Name
			}
			entityField.ColumnLength = columnMarker.Length
			entityField.UniqueColumn = columnMarker.Unique
		}

		if hasAnyMarker && entityField.IsTransient {
			err := fmt.Errorf("the '%s' transient field can only have shelf:transient marker", field.Name)
			errs = append(errs, marker.NewError(err, field.File.FullPath, position))
		}

		if hasAnyDateMarker {
			if TypeNameFromTypeWithoutImport(field.Type) != "Time" && importPath != "time" {
				err := fmt.Errorf("shelf:temporal, shelf:created-date and shelf:last-modified-date markers can be only applied to the time.Time type")
				errs = append(errs, marker.NewError(err, field.File.FullPath, position))
			}

			if isIdField {
				err := fmt.Errorf("id field '%s' cannot have shelf:temporal, shelf:created-date and shelf:last-modified-date markers", field.Name)
				errs = append(errs, marker.NewError(err, field.File.FullPath, position))
			}
		}

		if isIdField && hasTransientMarker {
			err := fmt.Errorf("the transient field '%s' cannot be id", field.Name)
			errs = append(errs, marker.NewError(err, field.File.FullPath, position))
		}

		if _, ok := columns[entityField.ColumnName]; ok {
			err := fmt.Errorf("there is already a column with name '%s'", field.Name)
			errs = append(errs, marker.NewError(err, field.File.FullPath, position))
		}

		columns[entityField.ColumnName] = true

		if isIdField {
			entityMetadata.IdField = entityField
		} else {
			entityMetadata.Fields = append(entityMetadata.Fields, entityField)
		}
	}
}

func ValidateFieldType(field marker.Field, typ marker.Type) bool {
	position := marker.Position{
		Line:   field.Position.Line,
		Column: field.Position.Column,
	}

	switch typedFieldType := typ.(type) {
	case *marker.AnonymousStructType:
		err := fmt.Errorf("the anonymous struct field is not supported for entities")
		errs = append(errs, marker.NewError(err, field.File.FullPath, position))
		return false
	case *marker.FunctionType:
		err := fmt.Errorf("the function field type is not supported for entities")
		errs = append(errs, marker.NewError(err, field.File.FullPath, position))
		return false
	case *marker.ChanType:
		err := fmt.Errorf("the chan field type is not supported for entities")
		errs = append(errs, marker.NewError(err, field.File.FullPath, position))
		return false
	case *marker.DictionaryType:
		err := fmt.Errorf("the map field type is not supported for entities")
		errs = append(errs, marker.NewError(err, field.File.FullPath, position))
		return false
	case *marker.ObjectType:
		return true
	case *marker.PointerType:
		return ValidateFieldType(field, typedFieldType.Typ)
	}

	return true
}
