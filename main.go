package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)


func main() {
	flag.Parse()
	templatesDir := os.Getenv("TEMPLATE")
	if templatesDir == "" {
		fmt.Println("You must define a $TEMPLATE environment variable.")
		return
	}
	if flag.Arg(0) == "" {
			template, err := os.Open(templatesDir)
			files, err := template.Readdir(-1)
			if  err != nil {
				fmt.Println("can' read template directory")
				return
			}
			for _ , file := range files {
				fmt.Println(file.Name())
			}
		return
	}
	templateDir := filepath.Join(templatesDir, flag.Arg(0))
	filepath.Walk(templateDir, display)
}

func display(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		fmt.Println("copy " + path )
	}
	return nil
}


