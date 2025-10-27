package configuration

import (
	"bytes"
	"encoding/base64"
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

	ensureCredentialFiles(t, params)

	basicTemplate.Execute(&templatedConfig, params)
	return templatedConfig.String()
}

func ensureCredentialFiles(t *testing.T, params map[string]string) {
	for key, value := range params {
		if strings.Contains(key, "Credentials") && strings.HasSuffix(value, ".json") {

			credentialsBase64 := readCredentialFileAsBase64(t, value)
			params[key] = credentialsBase64
		}
	}
}

func readCredentialFileAsBase64(t *testing.T, filename string) string {
	srcPath := "../test/fixtures/" + filename

	content, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("Could not read credential file %s: %v", srcPath, err)
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	t.Logf("Read credential file %s and encoded as base64 (%d characters)", srcPath, len(encoded))

	return encoded
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
