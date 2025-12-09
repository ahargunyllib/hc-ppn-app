package csv

import (
	"encoding/csv"
	"mime/multipart"
)

//go:generate mockgen -destination=mock/mock_csv.go -package=mock github.com/ahargunyllib/hc-ppn-app/apps/bot-service/pkg/csv CustomCSVInterface

type CustomCSVInterface interface {
	ParseFileHeader(fileHeader *multipart.FileHeader) ([][]string, error)
}

type CustomCSVStruct struct{}

var CSV = getCSV()

func getCSV() CustomCSVInterface {
	return &CustomCSVStruct{}
}

func (c *CustomCSVStruct) ParseFileHeader(fileHeader *multipart.FileHeader) ([][]string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
