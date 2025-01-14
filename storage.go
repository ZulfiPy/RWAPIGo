package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Storage[T any] struct {
	FileName string
}

func EnsureStorageFile[T any](storage *Storage[T], data T) error {
	dataPath := fmt.Sprintf("./%s", storage.FileName)

	if _, err := os.Stat(dataPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("creating file", storage.FileName)
			return storage.Save(data)
		}
		fmt.Println("Error accessing file:", err)
		return err
	}

	fmt.Printf("file exists %s, no action needed.\n", storage.FileName)
	return nil
}

func NewStorage[T any](fileName string) *Storage[T] {
	return &Storage[T]{FileName: fileName}
}

func (storage *Storage[T]) Save(data T) error {
	fileData, err := json.MarshalIndent(data, "", "    ")

	if err != nil {
		return err
	}

	return os.WriteFile(storage.FileName, fileData, 0644)
}

func (storage *Storage[T]) Load(data *T) error {
	fileData, err := os.ReadFile(storage.FileName)

	if err != nil {
		return err
	}

	return json.Unmarshal(fileData, data)
}
