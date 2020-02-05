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

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate <templatefile>",
	Short: "Validate your tosca template",
	Long: `Example:
dodas validate --template my_tosca_template.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf("Template: %v \n", string(templateFile))
			template, err := ioutil.ReadFile(templateFile)
			if err != nil {
				panic(err)
			}
			err = clientConf.Validate(template)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("Template: %v \n", string(args[0]))
			template, err := ioutil.ReadFile(args[0])
			if err != nil {
				panic(err)
			}
			err = clientConf.Validate(template)
			if err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")
	validateCmd.PersistentFlags().StringVar(&templateFile, "template", "", "Path to TOSCA template file")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
