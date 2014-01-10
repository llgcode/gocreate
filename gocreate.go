// Copyright 2012 The gocreate Authors. All rights reserved.
// created on 18/04/2012 by Laurent Le Goff

// Command line utility that create files from templates.
//     go get bitbucket.org/llg/gocreate
//     gocreate
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const DefaultTemplateDir = "bitbucket.org/llg/gocreate/templates"
const ConfigFileName = "config.json"

var help = flag.Bool("help", false, "Show help")
var override = flag.Bool("f", false, "Override existing files")

type Config struct {
	Doc        string
	LeftDelim  string
	RightDelim string
	Vars       map[string]interface{}
	Args       []*Arg
}

type Arg struct {
	Arg      string
	Name     string
	Doc      string
	Value    *string `json:"default"`
	Required bool
}

func readConfigFile(templateDir string) (c *Config) {
	configFilePath := filepath.Join(templateDir, ConfigFileName)
	content, err := ioutil.ReadFile(configFilePath)
	c = &Config{}
	if err == nil {
		err = json.Unmarshal(content, c)
		if err != nil {
			fmt.Println("Warning:", err)
		}
	}
	return
}

func showCommandHelp(cmd, templateDir string, c *Config) {
	fmt.Printf("Usage:\n  gocreate %s", cmd)
	for _, arg := range c.Args {
		if strings.HasPrefix(arg.Arg, "$") {
			fmt.Printf(" '%s'", arg.Name)
		}
	}
	fmt.Println("\n")
	fmt.Println("  " + c.Doc + "\n")
	fmt.Println("  -help: Show Command help")
	fmt.Println("  -f: Override existing files")
	for _, arg := range c.Args {
		fmt.Printf("  -%s=%s: '%s'. %s (required:%t)\n", arg.Arg, *arg.Value, arg.Name, arg.Doc, arg.Required)
	}
	fmt.Println("\n  Template Path:", templateDir)
}

func showHelp(templatesDir string) {
	fmt.Fprintln(os.Stderr, "Usage of gocreate:\n")
	fmt.Fprintln(os.Stderr, "  gocreate 'templateName'\n")
	flag.PrintDefaults()
	template, err := os.Open(templatesDir)
	files, err := template.Readdir(-1)
	if err != nil {
		fmt.Println("can't read template directory", templatesDir)
		return
	}
	defer template.Close()
	fmt.Println("\nTemplates Path:", templatesDir)
	fmt.Println("\nList of Templates: ")
	for _, file := range files {
		if file.IsDir() {
			c := readConfigFile(filepath.Join(templatesDir, file.Name()))
			fmt.Printf("%s: gocreate %s -help\n    %s\n", file.Name(), file.Name(), c.Doc)
		}
	}
}
func create(ctx *template.Template, sourceFilePath, destFilePath string, c *Config) {
	f, err := os.Open(destFilePath)
	alreadyExist := false
	if err == nil {
		alreadyExist = true
		f.Close()
	}
	if *override || !alreadyExist {
		if alreadyExist {
			fmt.Println("Override:", destFilePath)
		} else {
			fmt.Println("Create:", destFilePath)
		}
		os.IsExist(err)
		destFile, err := os.OpenFile(destFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}
		defer destFile.Close()
		tmpl, err := ctx.ParseFiles(sourceFilePath)
		if err != nil {
			panic(err)
		}
		tmplName := filepath.Base(sourceFilePath)
		err = tmpl.ExecuteTemplate(destFile, tmplName, c.Vars)
	} else {
		fmt.Println("Already exist don't override:", destFilePath)
	}
}

func createFromTemplateDir(ctx *template.Template, templateSourcePath, sourceFolder, destFolder string, c *Config) {
	templateDir, err := os.Open(sourceFolder)
	files, err := templateDir.Readdir(-1)
	if err != nil {
		fmt.Println("can't read template directory", sourceFolder)
		showHelp(filepath.Dir(templateSourcePath))
		return
	}
	defer templateDir.Close()
	for _, file := range files {
		sourceFilePath := filepath.Join(sourceFolder, file.Name())
		tmplName := filepath.Base(sourceFilePath)
		tmpl, err := template.New("filename").Delims(c.LeftDelim, c.RightDelim).Parse(tmplName)
		if err != nil {
			panic(err)
		}
		var buf bytes.Buffer
		err = tmpl.ExecuteTemplate(&buf, "filename", c.Vars)
		destFilePath := filepath.Join(destFolder, string(buf.Bytes()))
		if file.IsDir() {
			err = os.MkdirAll(destFilePath, 0766)
			if err != nil {
				panic(err)
			}
			createFromTemplateDir(ctx, templateSourcePath, sourceFilePath, destFilePath, c)
		} else if sourceFilePath != filepath.Join(templateSourcePath, ConfigFileName) {
			create(ctx, sourceFilePath, destFilePath, c)
		}
	}
}

func main() {
	templatesDirPath := os.Getenv("GOTEMPLATE")
	if templatesDirPath == "" {
		if list := filepath.SplitList(os.Getenv("GOPATH")); len(list) > 0 && list[0] != runtime.GOROOT() {
			templatesDirPath = filepath.Join(list[0], "src", DefaultTemplateDir)
		} else {
			templatesDirPath = filepath.Join(runtime.GOROOT(), "src", "pkg", DefaultTemplateDir)
		}
	}
	if len(os.Args) == 1 {
		showHelp(templatesDirPath)
		return
	}
	templateName := os.Args[1]
	if os.Args[1] == "-help" {
		flag.Parse()
		showHelp(templatesDirPath)
		return
	}

	templateDirPath := filepath.Join(templatesDirPath, templateName)
	c := readConfigFile(templateDirPath)
	for _, arg := range c.Args {
		if arg.Value == nil {
			arg.Value = new(string)
		}
		if !strings.HasPrefix(arg.Arg, "$") {
			arg.Value = flag.String(arg.Arg, *arg.Value, arg.Doc)
		}
	}
	os.Args = os.Args[1:]
	flag.Parse()
	if *help {
		showCommandHelp(templateName, templateDirPath, c)
		return
	}
	if c.Vars == nil {
		c.Vars = make(map[string]interface{})
	}
	for _, arg := range c.Args {
		var val string
		if strings.HasPrefix(arg.Arg, "$") {
			i, err := strconv.ParseInt(arg.Arg[1:], 10, 0)
			if err != nil {
				panic(err)
			}
			*arg.Value = flag.Arg(int(i))
		}
		val = *arg.Value
		if val == "" && arg.Required {
			fmt.Println("-"+arg.Name, " option is Required!!")
			showCommandHelp(templateName, templateDirPath, c)
			return
		}

		c.Vars[arg.Name] = val
	}
	c.Vars["now"] = time.Now()
	ctx := template.New("templates")

	templatesDir, err := os.Open(templatesDirPath)
	files, err := templatesDir.Readdir(-1)
	if err != nil {
		fmt.Println("can't read template directory", templatesDirPath)
		return
	}
	defer templatesDir.Close()
	for _, file := range files {
		if !file.IsDir() {
			ctx, err = ctx.ParseGlob(filepath.Join(templatesDirPath, file.Name()))
			if err != nil {
				panic(err)
			}
		}
	}
	ctx = ctx.Delims(c.LeftDelim, c.RightDelim)

	createFromTemplateDir(ctx, templateDirPath, templateDirPath, ".", c)

}
