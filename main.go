package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func parseEnv() map[string]interface{} {
	var result map[string]interface{}
	config := os.Getenv("APPCONFIG")
	json.Unmarshal([]byte(config), &result)

	return result
}

func createNewEnv(config map[string]interface{}) string {
	fileContents := "window._env_ = {\n"

	for key, value := range config {
		line := fmt.Sprintf("\t%s: \"%s\",\n", key, value)
		fileContents += line
	}

	fileContents += "};"

	return fileContents
}

func main() {
	envPath := os.Getenv("ENVPATH")
	os.Remove(envPath)
	os.Create(envPath)

	appConfig := parseEnv()
	fileContents := []byte(createNewEnv(appConfig))
	fileErr := ioutil.WriteFile(envPath, fileContents, 0)

	if fileErr != nil {
		panic(fileErr)
	}

	os.Exit(0)
}
