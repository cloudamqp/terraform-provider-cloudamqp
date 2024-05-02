package configuration

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"text/template"
)

func GetTemplatedConfig(t *testing.T, fileNames []string, params map[string]string) string {
	config := appendFiles(t, fileNames)

	var templatedConfig bytes.Buffer
	basicTemplate := template.Must(template.New("template").Parse(config))
	basicTemplate.Execute(&templatedConfig, params)
	return templatedConfig.String()
}

func appendFiles(t *testing.T, fileNames []string) string {
	configArray := make([]string, len(fileNames))
	for index, filename := range fileNames {
		file := fmt.Sprintf("%s.txt", filename)
		tempConfig, err := loadFile(file)
		if err != nil {
			t.Fatalf("failed to load configuration, err: %v", err)
		}
		configArray[index] = tempConfig
	}
	return strings.Join(configArray, "\n")
}

func loadFile(filename string) (string, error) {
	path := "../test/configurations"
	file, err := os.Open(fmt.Sprintf("%s/%s", path, filename))
	if err != nil {
		return "", err
	}
	defer file.Close()
	var rawFile bytes.Buffer
	_, err = io.Copy(&rawFile, file)
	if err != nil {
		return "", err
	}
	return rawFile.String(), nil
}
