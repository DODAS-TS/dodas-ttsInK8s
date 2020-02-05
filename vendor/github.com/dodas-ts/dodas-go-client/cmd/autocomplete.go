package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// autocompleteCmd represents the completion command
var autocompleteCmd = &cobra.Command{
	Use:   "autocomplete",
	Short: "Generate script for bash autocomplete",
	Long: `add the following line to ~/.bashrc: . <(dodas autocomplete)`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout);
		//rootCmd.GenZshCompletion(os.Stdout);
	},
}

// autocompleteZshCmd represents the completion command
var autocompleteZshCmd = &cobra.Command{
	Use:   "zsh-autocomplete",
	Short: "Generate script for zsh autocomplete",
	Long: `add the following line to ~/.bashrc: source <(dodas zsh-autocomplete)`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenZshCompletion(os.Stdout);
	},
}

func init() {
	rootCmd.AddCommand(autocompleteCmd)
	rootCmd.AddCommand(autocompleteZshCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completitionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completitionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
