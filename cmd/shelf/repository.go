package main

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/shelf"
	"strings"
)

var reservedRepositoryMethods = map[string]int{
	"Count":         2,
	"ExistsById":    2,
	"Delete":        2,
	"DeleteById":    2,
	"DeleteAll":     2,
	"DeleteAllById": 2,
	"Save":          2,
	"SaveAll":       2,
	"FindById":      2,
	"FindAll":       2,
	"FindAllById":   2,
}

type RepositoryMetadata struct {
	RepositoryName string
	EntityName     string
	InterfaceType  marker.InterfaceType
}

func ValidateRepositoryMarkers(interfaceType marker.InterfaceType) bool {
	markers := interfaceType.Markers

	if markers == nil {
		return true
	}

	isValid := true

	for name, values := range markers {
		if values != nil && len(values) > 1 {
			err := fmt.Errorf("the interface cannot be marked twice as '%s' marker", name)
			errs = append(errs, marker.NewError(err, interfaceType.File.FullPath, marker.Position{
				Line:   interfaceType.Position.Line,
				Column: interfaceType.Position.Column,
			}))
			isValid = false
		}
	}

	return isValid
}

func FindRepositories(interfaceTypes []marker.InterfaceType) {
	for _, interfaceType := range interfaceTypes {
		if !ValidateRepositoryMarkers(interfaceType) {
			continue
		}

		markerValues := interfaceType.Markers

		if markerValues == nil {
			continue
		}

		markers, ok := markerValues[shelf.MarkerRepository]

		if !ok {
			return
		}

		var err error
		entityName := ""
		repositoryName := strings.TrimSpace(interfaceType.Name)

		for _, candidateMarker := range markers {
			switch typedMarker := candidateMarker.(type) {
			case shelf.RepositoryMarker:
				entityName = strings.TrimSpace(typedMarker.Entity)
				repositoryNameValue := strings.TrimSpace(typedMarker.Name)

				if repositoryNameValue != "" {
					repositoryName = repositoryNameValue
				}

				if _, ok := repositoriesByName[repositoryName]; ok {
					err = fmt.Errorf("there is already a repository with name '%s'", repositoryName)
					errs = append(errs, marker.NewError(err, interfaceType.File.FullPath, marker.Position{
						Line:   interfaceType.Position.Line,
						Column: interfaceType.Position.Column,
					}))
					break
				}
			}
		}

		if err == nil {
			if _, ok := entitiesByName[entityName]; !ok {
				err = fmt.Errorf("entity with name '%s' does not exist, please use a valid entity name", entityName)
				errs = append(errs, marker.NewError(err, interfaceType.File.FullPath, marker.Position{
					Line:   interfaceType.Position.Line,
					Column: interfaceType.Position.Column,
				}))
				continue
			}

			metadata := RepositoryMetadata{
				RepositoryName: repositoryName,
				EntityName:     entityName,
				InterfaceType:  interfaceType,
			}

			fullInterfaceName := interfaceType.File.Package.Path + "#" + interfaceType.Name
			repositoryMetadataByInterfaceName[fullInterfaceName] = metadata
			repositoriesByName[repositoryName] = fullInterfaceName
		}

	}
}

func ValidateRepositoryMethods(methods []marker.Method) {
	for _, method := range methods {
		ValidateRepositoryMethodParameters(method)
		ValidateXMarkers(method)
	}
}

func ValidateXMarkers(method marker.Method) {
	markerValues := method.Markers

	if markerValues == nil {
		return
	}

	markers, ok := markerValues[shelf.MarkerQuery]

	if !ok {
		return
	}

	matched := false

	for _, candidateMarker := range markers {
		switch candidateMarker.(type) {
		case shelf.QueryMarker:

			if matched {
				err := fmt.Errorf("repository methods cannot be marked twice as '%s' marker", shelf.MarkerQuery)
				errs = append(errs, marker.NewError(err, method.File.FullPath, marker.Position{
					Line:   method.Position.Line,
					Column: method.Position.Column,
				}))
				break
			}

			matched = true
		}
	}
}

func ValidateRepositoryMethodParameters(method marker.Method) {

	if method.Parameters == nil || len(method.Parameters) < 1 {
		err := errors.New("repository methods must take in one parameter of type context.Context at least")
		errs = append(errs, marker.NewError(err, method.File.FullPath, marker.Position{
			Line:   method.Position.Line,
			Column: method.Position.Column,
		}))
	}

	for index, param := range method.Parameters {
		if index == 0 {
			name := GetFullNameFromType(param.Type)
			if "context.Context" != name {
				err := errors.New("the type of the first parameter must be context.Context for repositories")
				errs = append(errs, marker.NewError(err, method.File.FullPath, marker.Position{
					Line:   method.Position.Line,
					Column: method.Position.Column,
				}))
				break
			}
		}
	}
}
