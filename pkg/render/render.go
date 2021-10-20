package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/DaniellaFreese/go-course/pkg/config"
	"github.com/DaniellaFreese/go-course/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

//NewTemplate sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	tc := app.TemplateCache

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("no html template")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}

}

//CreateTemplateCache creates a etmpalte cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.go.html")

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			fmt.Println("error on creating template")
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.go.html")
		if err != nil {
			fmt.Println("error on getting matches")
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.go.html")
			if err != nil {
				fmt.Println("error on getting the layout")
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
