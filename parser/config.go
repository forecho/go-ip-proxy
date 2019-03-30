package parser

import (
	"encoding/json"
	"go-ip-proxy/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

type Configs struct {
	Configs []*Config `json:"config"`
}

type Config struct {
	Name          string `json:"name"`
	UrlFormat     string `json:"urlFormat"`
	UrlParameters string `json:"urlParameters"`
	Type          Type   `json:"collectType"`
	Charset       string `json:"charset"`
	ValueRuleMap  []struct {
		Name  string `json:"name"`
		Rule  string `json:"rule"`
		Value string `json:"value,omitempty"`
	} `json:"valueNameRuleMap"`
}

func NewConfig(fileName string) *Configs {
	configFile, err := os.Open(fileName)
	if err != nil {
		logger.Error("file does not exist")
	}

	defer configFile.Close()
	var configs Configs
	byteValue, _ := ioutil.ReadAll(configFile)
	jsonErr := json.Unmarshal(byteValue, &configs)
	if jsonErr != nil {
		logger.Error(fileName+" file json decode error", zap.Error(jsonErr))
	}

	return &configs
}
