package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func parseEnv() map[string]interface{} {
	var result map[string]interface{}
	config := os.Getenv("APPCONFIG")
	json.Unmarshal([]byte(config), &result)

	return result
}

func main() {
	err := filepath.Walk(os.Getenv("BUILD_DIR"), visit)

	if err != nil {
		panic(err)
	}

	os.Exit(0)
}

func visit(path string, fi os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	if !!fi.IsDir() {
		return nil
	}

	matched, err := filepath.Match("*.js", fi.Name())

	if err != nil {
		panic(err)
	}

	if matched {
		fmt.Printf("\nReplacing variables in: %s\n", fi.Name())
		read, err := ioutil.ReadFile(path)

		if err != nil {
			panic(err)
		}

		var newContents = string(read)

		for key, value := range parseEnv() {
			fmt.Printf("Replacing %s with value %s\n", key, value)
			r := regexp.MustCompile(fmt.Sprintf("(\"%s\")[:](\"[^\"]*\")", key))

			newContents = r.ReplaceAllString(newContents, fmt.Sprintf("\"%s\":\"%s\"", key, value))

			if err != nil {
				panic(err)
			}
		}

		err = ioutil.WriteFile(path, []byte(newContents), 0)

		if err != nil {
			panic(err)
		}
	}

	return nil
}
