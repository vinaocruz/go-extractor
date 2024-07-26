package service

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/cavaliergopher/grab/v3"
)

const (
	storage = "./storage/"
)

type DownloadManager interface {
	GetB3Files() []string
	ExtractZipFile(zipFile string, files *[]string, wg *sync.WaitGroup)
}

type LocalDownloadManager struct {
}

func NewDownloadManager() DownloadManager {
	return &LocalDownloadManager{}
}

func (dm *LocalDownloadManager) GetB3Files() []string {
	fmt.Println("downloading B3 files...")

	// TODO implementar crawler para buscar últimos arquivos
	// em https://www.b3.com.br/pt_br/market-data-e-indices/servicos-de-dados/market-data/cotacoes/cotacoes/
	remoteFiles := []string{
		"https://up2dataweb.blob.core.windows.net/raptor/01-07-2024_NEGOCIOSAVISTA.zip?sv=2019-12-12&ss=b&srt=o&spr=https&se=2024-07-07T19%3A32%3A07Z&sp=r&sig=LMi1yviURFR1fTvqJvmJ1NFZyhX%2BsPtE3wKEh74afDA%3D",
		"https://up2dataweb.blob.core.windows.net/raptor/28-06-2024_NEGOCIOSAVISTA.zip?sv=2019-12-12&ss=b&srt=o&spr=https&se=2024-07-07T19%3A32%3A57Z&sp=r&sig=J1H%2FA07fS3LQ3y3uyMzzqDHzRs0BWuz%2BqKxGE3gLTuY%3D",
	}

	// concurrency download
	responseChannel, err := grab.GetBatch(0, storage+"swap/", remoteFiles...)
	if err != nil {
		log.Fatal("Failed to download files: ", err)
	}

	var wg sync.WaitGroup

	var files []string
	for resp := range responseChannel {
		if err := resp.Err(); err != nil {
			log.Fatal("Failed to download file: ", err)
		}

		fmt.Printf("Completed download %s\n", resp.Request.URL())

		wg.Add(1)
		go dm.ExtractZipFile(resp.Filename, &files, &wg)
	}

	wg.Wait()
	fmt.Println("finished download!")

	return files
}

func (dm *LocalDownloadManager) ExtractZipFile(zipFile string, files *[]string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("extracting " + zipFile + "...")

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Fatal("Failed to open zip file: ", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			log.Fatal("Failed to open file: ", err)
		}
		defer rc.Close()

		if !file.FileInfo().IsDir() {
			filepath := fmt.Sprintf("%s%s", storage+"example/", file.Name)

			uncompressedFile, err := os.Create(filepath)
			if err != nil {
				log.Fatal("Failed to create file: ", err)
			}

			_, err = io.Copy(uncompressedFile, rc)
			if err != nil {
				log.Fatal("Failed to copy file: ", err)
			}

			*files = append(*files, filepath)
		}
	}
}
