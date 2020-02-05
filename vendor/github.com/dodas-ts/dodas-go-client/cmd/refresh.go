package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Not implemented yet",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not implemented yet")

		// tok := "CHANGEME"

		// r := RefreshRequest{
		// 	Endpoint:     "https://dodas-iam.cloud.cnaf.infn.it/token",
		// 	ClientID:     "CHANGEME",
		// 	ClientSecret: "CHANGEME",
		// 	AccessToken:  tok,
		// }
		// token, err := GetRefreshToken(r)
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println(token)

		// r.RefreshToken = token

		// t, err := GetAccessToken(r)
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Printf("Access Token: %s", t)

	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// refreshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// refreshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
