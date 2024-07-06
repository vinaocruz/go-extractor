/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	importShortDesc = "A brief description of your command"
	importLongDesc  = "A longer description that spans multiple lines and likely contains examples"
)

// dataImportCmd represents the dataImport command
var dataImportCmd = &cobra.Command{
	Use:   "dataImport",
	Short: importShortDesc,
	Long:  importLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		execute()
	},
}

func execute() {
	downloadData()
	extractData()
	insertData()
}

func downloadData() {
	fmt.Println("downloading...")
}

func extractData() {
	fmt.Println("extracting...")
}

func insertData() {
	fmt.Println("inserting...")
}

func init() {
	rootCmd.AddCommand(dataImportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dataImportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dataImportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
