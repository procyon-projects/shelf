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

var paths []string
var outputPath string
var options []string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go files by processing markers",
	Long:  `The generate command helps your code generation process by processing markers`,
	Run: func(cmd *cobra.Command, args []string) {
		packages, err := marker.LoadPackages(paths...)

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
			return
		}

		err = ProcessMarkers(collector, packages)

		if err != nil {
			PrintError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringSliceVarP(&paths, "path", "p", paths, "path(s) separated by comma")
	err := generateCmd.MarkFlagRequired("path")

	if err != nil {
		panic(err)
	}

	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "", "output path")
	err = generateCmd.MarkFlagRequired("output")

	if err != nil {
		panic(err)
	}

	generateCmd.Flags().StringSliceVarP(&options, "args", "a", options, "extra arguments for marker processors (key-value separated by comma)")
}
