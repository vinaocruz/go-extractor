/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/vinaocruz/go-extractor/src/model"
	"github.com/vinaocruz/go-extractor/src/repository"
	"github.com/vinaocruz/go-extractor/src/service"
)

type DataImportCmd struct {
	Repo        repository.NegociationRepository
	Service     service.FileManager
	DownService service.DownloadManager
}

const (
	importShortDesc = "A brief description of your command"
	importLongDesc  = "A longer description that spans multiple lines and likely contains examples"
)

func NewDataImportCmd() *cobra.Command {
	dataImport := &DataImportCmd{
		Repo:        repository.NewNegociationRepository(),
		Service:     service.NewLocalFileManager(),
		DownService: service.NewDownloadManager(),
	}

	var dataImportCmd = &cobra.Command{
		Use:   "dataImport",
		Short: importShortDesc,
		Long:  importLongDesc,
		Run: func(cmd *cobra.Command, args []string) {
			files := []string{
				"storage/example/27-06-2024_NEGOCIOSAVISTA.zip",
				"storage/example/28-06-2024_NEGOCIOSAVISTA.zip",
				"storage/example/01-07-2024_NEGOCIOSAVISTA.zip",
				"storage/example/02-07-2024_NEGOCIOSAVISTA.zip",
				"storage/example/03-07-2024_NEGOCIOSAVISTA.zip",
				"storage/example/04-07-2024_NEGOCIOSAVISTA.zip",
				"storage/example/05-07-2024_NEGOCIOSAVISTA.zip",
			}

			dataImport.execute(files)
		},
	}

	return dataImportCmd
}

func (dm *DataImportCmd) execute(zipFiles []string) {
	files := dm.extract(zipFiles)

	var wg sync.WaitGroup

	batchCh := make(chan []model.Negociation, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			dm.Service.ReadFile(filename, batchCh)
		}(file)
	}

	go func() {
		for lines := range batchCh {
			start := time.Now()
			dm.Repo.BulkImport(lines)
			elapsed := time.Since(start)
			log.Printf("Bulk imported in %s", elapsed)
		}
	}()

	wg.Wait()
	close(batchCh)

	dm.Repo.SetupIndex()

	fmt.Println("Import finished!")
}

func (dm *DataImportCmd) extract(zipFiles []string) []string {
	var wg sync.WaitGroup
	var files []string

	for _, zipFile := range zipFiles {
		wg.Add(1)
		go dm.DownService.ExtractZipFile(zipFile, &files, &wg)
	}

	wg.Wait()

	return files
}

func init() {
	rootCmd.AddCommand(
		NewDataImportCmd(),
	)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dataImportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dataImportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
