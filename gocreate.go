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
	"text/template"
)

const DefaultTemplateDir = "bitbucket.org/llg/gocreate/template"
const ConfigFileName = "config.json"

var help = flag.Bool("help", false, "Show help")
var override = flag.Bool("f", false, "Override existing file")

type Config struct {
	Doc  string
	LeftDelim string
	RightDelim string
	Vars map[string]string
	Args []*Arg
}

type Arg struct {
	Arg   string
	Name  string
	Doc   string
	Value *string `json:"default"`
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
	fmt.Fprintf(os.Stderr, "Usage of gocreate %s:\n", cmd)
	fmt.Println("   ", c.Doc)
    flag.PrintDefaults()
	fmt.Println("    Template Path:", templateDir)
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
		c := readConfigFile(filepath.Join(templatesDir, file.Name()))
		fmt.Printf("%s: gocreate %s -help\n    %s\n", file.Name(), file.Name(), c.Doc)
	}
}
func create(sourceFile, destFolder string, c *Config) {
	sourceFileName := filepath.Base(sourceFile)
	tmpl, err := template.New("filename").Delims(c.LeftDelim, c.RightDelim).Parse(sourceFileName)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "filename", c.Vars)
	destFilePath := filepath.Join(destFolder, string(buf.Bytes()))
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
		destFile, err := os.OpenFile(destFilePath, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		defer destFile.Close()
		tmpl, err = template.New("file").Delims(c.LeftDelim, c.RightDelim).ParseFiles(sourceFile)
		if err != nil {
			panic(err)
		}
		err = tmpl.ExecuteTemplate(destFile, sourceFileName, c.Vars)
	} else {
		fmt.Println("Already exist don't override:", destFilePath)
	}
}

func createFromTemplateDir(templateSourcePath, sourceFolder, destFolder string, c *Config) {
	templateDir, err := os.Open(sourceFolder)
	files, err := templateDir.Readdir(-1)
	if err != nil {
		fmt.Println("can't read template directory", sourceFolder)
		showHelp(filepath.Dir(templateSourcePath))
		return
	}
	defer templateDir.Close()
	for _, file := range files {
		sourcePath := filepath.Join(sourceFolder, file.Name())
		if file.IsDir() {
			tmpl, err := template.New("filename").Delims(c.LeftDelim, c.RightDelim).Parse(file.Name())
			if err != nil {
				panic(err)
			}
			var buf bytes.Buffer
			err = tmpl.ExecuteTemplate(&buf, "filename", c.Vars)
			newDestFolder := filepath.Join(destFolder, string(buf.Bytes()))

			err = os.MkdirAll(newDestFolder, 0666)
			if err != nil {
				panic(err)
			}
			createFromTemplateDir(templateSourcePath, sourcePath, newDestFolder, c)
		} else if sourcePath != filepath.Join(templateSourcePath, ConfigFileName) {
			create(sourcePath, destFolder, c)
		}
	}
}

func main() {

	templatesDir := os.Getenv("GOTEMPLATE")
	if templatesDir == "" {
		if list := filepath.SplitList(os.Getenv("GOPATH")); len(list) > 0 && list[0] != runtime.GOROOT() {
			templatesDir = filepath.Join(list[0], "src", DefaultTemplateDir)
		} else {
			templatesDir = filepath.Join(runtime.GOROOT(), "src", "pkg", DefaultTemplateDir)
		}
	}
	if len(os.Args) == 1 {
		showHelp(templatesDir)
		return
	}
	templateName := os.Args[1]
	if os.Args[1] == "-help" {
		flag.Parse()
		showHelp(templatesDir)
		return
	}

	templateDirPath := filepath.Join(templatesDir, templateName)
	c := readConfigFile(templateDirPath)
	for _, arg := range c.Args {
		def := ""
		if arg.Value != nil {
			def = *arg.Value
		}
		arg.Value = flag.String(arg.Arg, def, arg.Doc)
	}
	os.Args = os.Args[1:]
	flag.Parse()
	if *help {
		showCommandHelp(templateName, templateDirPath, c)
		return
	}
	if c.Vars == nil {
		c.Vars = make(map[string]string)
	}
	for _, arg := range c.Args {
		val := *arg.Value
		if val == "" && arg.Required {
			fmt.Println("-" + arg.Name, " option is Required!!")
			showCommandHelp(templateName, templateDirPath, c)
			return
		}
		c.Vars[arg.Name] = *arg.Value
	}
	createFromTemplateDir(templateDirPath, templateDirPath, ".", c)

}
