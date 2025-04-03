package main

import (
	"blindspot/internal/challenges"
	"blindspot/internal/cli"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var (
	difficulty int
	category   string
)

var rootCmd = &cobra.Command{
	Use:   "blindspot",
	Short: "blindspot allows access to blindspot from the command line",
	Long: `blindspot is a terminal game made to help you learn how to identify
insecure code practices and train yourself!`,
	Run: func(cmd *cobra.Command, args []string) {
		config := cli.NewConfig()
		config.SetDifficulty(difficulty)
		config.SetFlagChanged("difficulty", cmd.Flags().Changed("difficulty"))
		config.SetCategory(category)
		config.SetFlagChanged("category", cmd.Flags().Changed("category"))

		runner := cli.NewRunner(config)
		runner.Run()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	categoriesHelp := "Set game vulnerability category"

	categories, err := challenges.GetChallengesCategories()
	if err == nil && len(categories) > 0 {
		categoriesStr := strings.Join(categories, ",\n ")
		categoriesHelp = fmt.Sprintf("Set game vulnerability category.\nAvailable categories:\n %s.", categoriesStr)
	}

	rootCmd.Flags().IntVarP(&difficulty, "difficulty", "d", 0,
		"Set game difficulty (0=beginner, 1=intermediate, 2=advanced)")
	rootCmd.Flags().StringVarP(&category, "category", "c", "",
		categoriesHelp)
}

func main() {
	if err := Execute(); err != nil {
		log.Fatal("Error executing command:", err)
	}
}
