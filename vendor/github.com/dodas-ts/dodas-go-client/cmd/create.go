// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var decodeFields = map[string]string{
	"ID":            "id",
	"Type":          "type",
	"Username":      "username",
	"Password":      "password",
	"Token":         "token",
	"Host":          "host",
	"Tenant":        "tenant",
	"AuthURL":       "auth_url",
	"AuthVersion":   "auth_version",
	"Domain":        "domain",
	"ServiceRegion": "service_region",
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <Template1> ... <TemplateN>",
	Args:  cobra.MinimumNArgs(1),
	Short: "Create a cluster from a TOSCA template",
	Long: `
`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, temp := range args {
			templateFile = temp
			fmt.Printf("Template: %v \n", string(templateFile))
			template, err := ioutil.ReadFile(templateFile)
			if err != nil {
				panic(err)
			}

			err = clientConf.Validate(template)
			if err != nil {
				panic(err)
			}

			_, err = clientConf.CreateInf(template)
			if err != nil {
				panic(err)
			}

		}
	},
}

func init() {

	rootCmd.AddCommand(createCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose local command")
}
