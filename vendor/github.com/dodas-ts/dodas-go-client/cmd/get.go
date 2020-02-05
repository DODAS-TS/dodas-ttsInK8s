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
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var debug bool

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.MinimumNArgs(1),
	Short: "Wrapper command for get operations",
	Long: `Wrapper command for get operations.
dodas get -h for possible commands`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// outputCmd represents the output command
var outputCmd = &cobra.Command{
	Use:   "output <infID>",
	Args:  cobra.MinimumNArgs(1),
	Short: "Get deployment output",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")

		outputs, err := clientConf.GetInfOutputs(string(args[0]))
		if err != nil {
			panic(err)
		}

		for key, value := range outputs {
			fmt.Println(key, value)
		}

	},
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status <infID>",
	Args:  cobra.MinimumNArgs(1),
	Short: "Get deployment status",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")
		authHeader := clientConf.PrepareAuthHeaders()

		request := Request{
			URL:         string(clientConf.Im.Host) + "/" + string(args[0]) + "/contmsg",
			RequestType: "GET",
			Headers: map[string]string{
				"Authorization": authHeader,
				"Content-Type":  "application/json",
			},
		}

		body, statusCode, err := MakeRequest(request)
		if err != nil {
			panic(err)
		}

		fmt.Print("Deployment status:\n")

		if statusCode == 200 {
			fmt.Println(string(body))
		} else {
			fmt.Println("ERROR:\n", string(body))
			return
		}

	},
}

// vmstatusCmd represents the status command
var vmstatusCmd = &cobra.Command{
	Use:   "vm <infID> <vmID>",
	Args:  cobra.MinimumNArgs(2),
	Short: "Get VM deployment logs",
	Long: `
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vmstatus called")

		status, err := clientConf.GetInfVMStates(string(args[0]), string(args[1]))
		if err != nil {
			panic(err)
		}

		if debug {
			scanner := bufio.NewScanner(strings.NewReader(status))

			reExclude := regexp.MustCompile(`contmsg`)
			reFinished := regexp.MustCompile(`^Task.*finished`)
			reStarted := regexp.MustCompile(`^Launch task:`)

			var line string

			for scanner.Scan() {

				line = scanner.Text()

				if (reFinished.Match([]byte(line)) || reStarted.Match([]byte(line))) && !reExclude.Match([]byte(line)) {
					fmt.Println(line)
				}
			}
		} else {
			fmt.Println(status)
		}

	},
}

// vmCmd represents the vm command
var vmCmd = &cobra.Command{
	Use:   "vm <infID> <vmID>",
	Args:  cobra.MinimumNArgs(2),
	Short: "Get VM details",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vm called")
		listVMs, err := GetVMs(string(args[0]))
		if err != nil {
			panic(err)
		}

		vmN := string(args[1])
		vmID, err := strconv.Atoi(vmN)
		if err != nil {
			panic(err)
		}
		vm := listVMs[vmID]

		authHeader := clientConf.PrepareAuthHeaders()
		request := Request{
			URL:         vm,
			RequestType: "GET",
			Headers: map[string]string{
				"Authorization": authHeader,
				"Content-Type":  "application/json",
			},
		}

		body, statusCode, err := MakeRequest(request)
		if err != nil {
			panic(err)
		}

		if statusCode == 200 {
			fmt.Println(string(body))
		} else {
			panic(fmt.Errorf("Server response code %d: %s", statusCode, body))
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.AddCommand(outputCmd)
	getCmd.AddCommand(statusCmd)
	statusCmd.AddCommand(vmstatusCmd)

	getCmd.AddCommand(vmCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	vmstatusCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable verbose status logging")
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
