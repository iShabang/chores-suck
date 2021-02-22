package views

import (
	"html/template"
	"io/ioutil"
	"strings"
)

var (
	templates *template.Template
	dashboard *template.Template
)

// LoadTemplates grabs all templates from the templates directory and stores them in memory
func LoadTemplates() error {
	dir := ""
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	var fileNames []string
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".html") {
			fileNames = append(fileNames, dir+name)
		}
	}

	templates, err = template.ParseFiles(fileNames...)
	if err != nil {
		return err
	}
	dashboard = templates.Lookup("dashboard.html")

	return nil
}
