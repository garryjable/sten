package engine

import (
	"encoding/json"
	"os"
)

// type Dictionary map[string]string

// func Load(path string) (Dictionary, error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	var dict Dictionary
// 	decoder := json.NewDecoder(file)
// 	if err := decoder.Decode(&dict); err != nil {
// 		return nil, err
// 	}
// 	return dict, nil
// }

type Dictionary map[string]string

func LoadDictionary(path string) (Dictionary, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var dict Dictionary
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&dict)
	return dict, err
}

func (d Dictionary) Lookup(stroke string) (string, bool) {
	result, ok := d[stroke]
	return result, ok
}
