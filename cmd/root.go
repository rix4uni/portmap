package cmd

import (
	"fmt"
	"os"

	"github.com/rix4uni/portmap/banner"
	"github.com/spf13/cobra"
)

var silent bool // Define silent as a global variable

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "portmap",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Check if the version flag is set
		if v, _ := cmd.Flags().GetBool("version"); v {
			banner.PrintVersion() // Print the version and exit
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Execute Cobra command processing
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Define flags
	rootCmd.Flags().BoolP("version", "v", false, "Print the version of the tool and exit.")
	rootCmd.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "Suppress banner output")

	// Print banner only if silent mode is not enabled
	cobra.OnInitialize(func() {
		if !silent {
			banner.PrintBanner()
		}
	})
}
