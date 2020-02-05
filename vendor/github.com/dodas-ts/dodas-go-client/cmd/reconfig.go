package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// reconfigCmd restart cluster configuration
var reconfigCmd = &cobra.Command{
	Use:   "reconfig <infID>",
	Args:  cobra.MinimumNArgs(1),
	Short: "restart cluster configuration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")

		client := &http.Client{
			Timeout: 300 * time.Second,
		}
		fmt.Println("Submitting request to  : ", clientConf.Im.Host)

		req, err := http.NewRequest("PUT", string(clientConf.Im.Host)+"/"+string(args[0])+"/reconfigure", nil)

		req.Header.Set("Content-Type", "application/json")

		authHeader := clientConf.PrepareAuthHeaders()

		req.Header.Set("Authorization", authHeader)

		var request []string
		for name, headers := range req.Header {
			name = strings.ToLower(name)
			for _, h := range headers {
				request = append(request, fmt.Sprintf("%v: %v", name, h))
			}
		}

		request = append(request, fmt.Sprint("\n"))
		//fmt.Printf(strings.Join(request, "\n"))

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode == 200 {
			fmt.Println("Command received correctly. Reconfiguration of the cluster will start soon.")
		} else {
			fmt.Println("ERROR:\n", string(body))
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(reconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
