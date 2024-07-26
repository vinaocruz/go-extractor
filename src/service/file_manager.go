package service

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/vinaocruz/go-extractor/src/model"
)

type FileManager interface {
	ReadFile(file string, batchCh chan<- model.Negociation, wg *sync.WaitGroup)
}

type LocalFileManager struct {
}

func NewLocalFileManager() FileManager {
	return &LocalFileManager{}
}

func (fm *LocalFileManager) ReadFile(file string, batchCh chan<- model.Negociation, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("reading " + file + "...")

	f, err := os.Open(file)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)

	ignoreFirstLine := true
	for scan.Scan() {
		if ignoreFirstLine {
			ignoreFirstLine = false
			continue
		}

		batchCh <- fm.parseLine(scan.Text())
	}

	if err := scan.Err(); err != nil {
		log.Fatal("Failed to read file: ", err)
	}
}

func (fm *LocalFileManager) parseLine(line string) model.Negociation {
	part := strings.Split(line, ";")

	price, err := strconv.ParseFloat(strings.Replace(part[3], ",", ".", -1), 32)
	if err != nil {
		log.Fatal(err)
	}

	quantity, err := strconv.Atoi(part[4])
	if err != nil {
		log.Fatal(err)
	}

	return model.Negociation{
		ClosedAt:      part[5],
		TransactionAt: part[8],
		TicketCode:    part[1],
		Price:         float32(price),
		Quantity:      quantity,
	}
}
