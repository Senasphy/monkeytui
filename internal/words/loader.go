package words

import (
	"encoding/json"

	assetswords "monkeytui/assets/words"
)

const DefaultPath = "english_1000.json"

type List struct {
	Language string   `json:"language"`
	Words    []string `json:"words"`
}

func Load(path string) ([]string, error) {
	if path == "" {
		path = DefaultPath
	}

	jsonFile, err := assetswords.FS.ReadFile(path)
	if err != nil {
		return nil, err
	}

	wordList := List{}
	if err := json.Unmarshal(jsonFile, &wordList); err != nil {
		return nil, err
	}

	return wordList.Words, nil
}
