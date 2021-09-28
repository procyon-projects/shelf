/*
Copyright © 2021 Shelf Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/shelf"
)

var (
	errs         []error
	entities     []marker.StructType
	repositories []marker.InterfaceType
	loaders      []marker.InterfaceType
)

// Register your marker definitions.
func RegisterDefinitions(registry *marker.Registry) error {
	markers := []struct {
		Name   string
		Level  marker.TargetLevel
		Output interface{}
	}{
		{Name: shelf.MarkerEntity, Level: marker.StructTypeLevel, Output: &shelf.EntityMarker{}},
		{Name: shelf.MarkerTable, Level: marker.StructTypeLevel, Output: &shelf.TableMarker{}},

		{Name: shelf.MarkerId, Level: marker.FieldLevel, Output: &shelf.IdMarker{}},
		{Name: shelf.MarkerGeneratedValue, Level: marker.FieldLevel, Output: &shelf.GeneratedValueMarker{}},
		{Name: shelf.MarkerColumn, Level: marker.FieldLevel, Output: &shelf.ColumnMarker{}},
		{Name: shelf.MarkerLob, Level: marker.FieldLevel, Output: &shelf.LobMarker{}},
		{Name: shelf.MarkerTransient, Level: marker.FieldLevel, Output: &shelf.TransientMarker{}},
		{Name: shelf.MarkerEnumerated, Level: marker.FieldLevel, Output: &shelf.EnumeratedMarker{}},

		{Name: shelf.MarkerRepository, Level: marker.InterfaceTypeLevel, Output: &shelf.RepositoryMarker{}},
		{Name: shelf.MarkerQuery, Level: marker.InterfaceMethodLevel, Output: &shelf.QueryMarker{}},
	}

	for _, m := range markers {
		err := registry.Register(m.Name, PkgId, m.Level, m.Output)
		if err != nil {
			return err
		}
	}

	return nil
}

// Process your markers.
func ProcessMarkers(collector *marker.Collector, pkgs []*marker.Package) error {
	marker.EachFile(collector, pkgs, func(file *marker.File, err error) {
		FindRepositories(file.InterfaceTypes)
	})
	return marker.NewErrorList(errs)
}
