package main

import (
	"log"

	"blindspot/internal/cli"

	"github.com/spf13/cobra"
)

var (
	difficulty int
	directMode bool
)

// base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "blindspot",
	Short: "blindspot allows access to blindspot from the command line",
	Long: `blindspot is a terminal game made to help you learn how to identify
insecure code practices and train yourself!`,
	Run: func(cmd *cobra.Command, args []string) {
		config := cli.NewConfig()
		config.SetDifficulty(difficulty)
		config.SetFlagChanged("difficulty", cmd.Flags().Changed("difficulty"))

		runner := cli.NewRunner(config)
		runner.Run()
	},
}

// adds all child commands to the root command and sets flags appropriately
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().IntVarP(&difficulty, "difficulty", "d", 0,
		"Set game difficulty (0=beginner, 1=intermediate, 2=advanced)")
	rootCmd.Flags().BoolVarP(&directMode, "direct", "r", false,
		"Run directly in challenge mode")
}

func main() {
	if err := Execute(); err != nil {
		log.Fatal("Error executing command:", err)
	}
}
