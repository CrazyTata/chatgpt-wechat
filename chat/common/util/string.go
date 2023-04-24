package util

import (
	"encoding/csv"
	"github.com/360EntSecGroup-Skylar/excelize"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"os"
)

func GetExcelDataByPath(filename string) ([][]string, error) {
	fileHandler, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	rows := fileHandler.GetRows("Sheet1")
	return rows, err
}

func GetCSVDataByPath(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 将CSV文件的字节数据转换为io.Reader接口
	f := transform.NewReader(file, simplifiedchinese.GBK.NewDecoder())

	reader := csv.NewReader(f)
	result := make([][]string, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, record)
	}
	return result, nil
}
