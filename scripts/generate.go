package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func main() {
	target := os.Args[1] + ".go"
	directory := os.Args[2]

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}

	content := "// Code generated by go generate; DO NOT EDIT.\n"
	content += "package main\n\n"

	for _, script := range files {
		fileName := script.Name()
		if !script.IsDir() && path.Ext(fileName) == ".groovy" {
			dot := strings.LastIndex(fileName, ".")
			name := fileName[0:dot]

			bytes, err := ioutil.ReadFile(path.Join(directory, fileName))
			if err != nil {
				panic(err)
			}

			constName := strings.Replace(name, "-", "_", -1)
			constName = strings.ToUpper(constName)

			content += "const " + constName + " = `" + string(bytes) + "`\n"

		}
	}

	err = ioutil.WriteFile(target, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}
