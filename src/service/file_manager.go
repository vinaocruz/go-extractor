package service

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/vinaocruz/go-extractor/src/model"
)

const (
	batchSize = 500000
)

type FileManager interface {
	ReadFile(file string, batchCh chan<- []model.Negociation)
}

type LocalFileManager struct {
}

func NewLocalFileManager() *LocalFileManager {
	return &LocalFileManager{}
}

func (fm *LocalFileManager) ReadFile(file string, batchCh chan<- []model.Negociation) {
	fmt.Println("reading " + file + "...")

	f, err := os.Open(file)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)

	var lines []model.Negociation

	ignoreFirstLine := true
	for scan.Scan() {
		if ignoreFirstLine {
			ignoreFirstLine = false
			continue
		}

		lines = append(lines, fm.parseLine(scan.Text()))
		if len(lines) >= batchSize {
			batchCh <- lines
			lines = []model.Negociation{}
		}
	}

	if len(lines) > 0 {
		batchCh <- lines
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
