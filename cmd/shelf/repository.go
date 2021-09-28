package main

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/shelf"
)

func FindRepositories(interfaceTypes []marker.InterfaceType) {
	for _, interfaceType := range interfaceTypes {
		markerValues := interfaceType.Markers

		if markerValues == nil {
			continue
		}

		markers, ok := markerValues[shelf.MarkerRepository]

		if !ok {
			return
		}

		matched := false

		for _, candidateMarker := range markers {
			switch candidateMarker.(type) {
			case shelf.RepositoryMarker:
				ValidateRepositoryMethods(interfaceType.Methods)

				if matched {
					err := fmt.Errorf("an interface cannot be marked twice as '%s' marker", shelf.MarkerRepository)
					errs = append(errs, marker.NewError(err, interfaceType.File.FullPath, marker.Position{
						Line:   interfaceType.Position.Line,
						Column: interfaceType.Position.Column,
					}))

					break
				}

				repositories = append(repositories, interfaceType)
				matched = true
			}
		}

	}
}

func ValidateRepositoryMethods(methods []marker.Method) {
	for _, method := range methods {
		ValidateRepositoryMethodParameters(method)
		ValidateRepositoryMarkers(method)
	}
}

func ValidateRepositoryMarkers(method marker.Method) {
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
