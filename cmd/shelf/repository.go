package main

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/shelf"
	"strings"
)

type RepositoryMetadata struct {
	RepositoryName string
	EntityName     string
	InterfaceType  marker.InterfaceType
}

func FindRepositoryTypes(interfaceTypes []marker.InterfaceType) {
	for _, interfaceType := range interfaceTypes {

		var repositoryMarker shelf.RepositoryMarker

		if value, ok := InterfaceMustBeMarkedOnceAtMost(interfaceType, shelf.MarkerRepository); ok {
			repositoryMarker = value.(shelf.RepositoryMarker)
		} else {
			continue
		}

		metadata := &RepositoryMetadata{
			RepositoryName: strings.TrimSpace(interfaceType.Name),
			EntityName:     strings.TrimSpace(repositoryMarker.Entity),
			InterfaceType:  interfaceType,
		}

		if strings.TrimSpace(repositoryMarker.Name) != "" {
			metadata.RepositoryName = repositoryMarker.Name
		}

		if _, ok := repositoriesByName[metadata.RepositoryName]; ok {
			err := fmt.Errorf("there is already a repository with name '%s'", metadata.RepositoryName)
			errs = append(errs, marker.NewError(err, interfaceType.File.FullPath, interfaceType.Position))
			continue
		}

		if _, ok := entitiesByName[metadata.EntityName]; !ok {
			err := fmt.Errorf("the entity with name '%s' does not exist, please use a valid entity name", metadata.EntityName)
			errs = append(errs, marker.NewError(err, interfaceType.File.FullPath, interfaceType.Position))
			continue
		}

		fullInterfaceName := FullInterfaceName(interfaceType)
		repositoryMetadataByInterfaceName[fullInterfaceName] = metadata
		repositoriesByName[metadata.RepositoryName] = fullInterfaceName
	}
}

func ProcessRepositories() {
	for _, _ = range repositoryMetadataByInterfaceName {

	}
}
