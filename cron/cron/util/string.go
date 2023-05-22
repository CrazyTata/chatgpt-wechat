package util

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

func GetFileSuffix(url string) string {
	arr := strings.Split(url, ".")
	return arr[len(arr)-1]
}

func DeleteFileIfExist(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// 文件不存在，直接返回
		return nil
	}
	if err != nil {
		// 其他错误，返回错误信息
		return err
	}

	if fileInfo.IsDir() {
		// 是个文件夹，返回错误
		return fmt.Errorf("%s is a directory", filePath)
	}

	// 删除文件
	if err := os.Remove(filePath); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func RenameFileIfExists(filePath string) (string, error) {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// 文件不存在，直接返回原文件路径
		return filePath, nil
	}
	if err != nil {
		// 其他错误，返回错误信息
		return "", err
	}

	if fileInfo.IsDir() {
		// 是个文件夹，返回错误
		return "", fmt.Errorf("%s is a directory", filePath)
	}

	// 文件存在，重命名
	now := time.Now().Format("20060102150405") // 生成当前时间戳
	newFilePath := fmt.Sprintf("%s_%s", filePath, now)
	if err := os.Rename(filePath, newFilePath); err != nil {
		log.Println(err)
		return "", err
	}

	return newFilePath, nil
}

func InArray(slice []string, element string) bool {
	i := sort.SearchStrings(slice, element)
	return i < len(slice) && slice[i] == element
}
