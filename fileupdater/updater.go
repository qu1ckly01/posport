package fileupdater

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var keyMap = map[string]string{
	"Seria":     "Series",
	"Nomer":     "Number",
	"Umia":      "FirstName",
	"Famlia":    "LastName",
	"Otchestvo": "MiddleName",
}

func UpdateJSONFiles(srcDir, dstDir string) {
	for {
		files, err := os.ReadDir(srcDir)
		if err != nil {
			log.Println("Ошибка чтения папки:", err)
			return
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
				path := filepath.Join(srcDir, file.Name())
				data, err := os.ReadFile(path)
				if err != nil {
					log.Println("Ошибка чтения файла:", err)
					continue
				}

				var jsonObj map[string]interface{}
				err = json.Unmarshal(data, &jsonObj)
				if err != nil {
					log.Println("Ошибка парсинга JSON:", err)
					continue
				}

				newJson := make(map[string]interface{})
				for k, v := range jsonObj {
					if newKey, ok := keyMap[k]; ok {
						newJson[newKey] = v
					} else {
						newJson[k] = v
					}
				}

				newData, _ := json.MarshalIndent(newJson, "", "  ")
				dstPath := filepath.Join(dstDir, file.Name())
				_ = os.WriteFile(dstPath, newData, 0644)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
