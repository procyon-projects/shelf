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
	"github.com/spf13/cobra"
	"log"
)

var (
	validatePaths    []string
	validateArgs     []string
	validationErrors []error
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate markers' syntax and arguments",
	Long:  `The validate command helps you validate markers' syntax and arguments'`,
	Run: func(cmd *cobra.Command, args []string) {
		packages, err := marker.LoadPackages(validatePaths...)

		if err != nil {
			log.Errorln(err)
			return
		}

		registry := marker.NewRegistry()
		err = RegisterDefinitions(registry)

		if err != nil {
			log.Errorln(err)
			return
		}

		collector := marker.NewCollector(registry)
		err = ValidateMarkers(collector, packages)

		if err != nil {
			PrintError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringSliceVarP(&validatePaths, "path", "p", validatePaths, "path(s) separated by comma")
	err := validateCmd.MarkFlagRequired("path")

	if err != nil {
		panic(err)
	}

	validateCmd.Flags().StringSliceVarP(&validateArgs, "args", "a", validateArgs, "extra arguments for marker processors (key-value separated by comma)")
}
