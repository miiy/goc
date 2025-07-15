package template

import (
	"github.com/gin-contrib/multitemplate"
	"html/template"
	"io/fs"
	"path"
)

type Template struct {
	Render multitemplate.Renderer
	FuncMap map[string]interface{}
}

func (r *Template) AddFunc(name string, i interface{}) {
	r.FuncMap[name] = i
}

// AddTemplate
// Example:
// filesMap = map[string][]string{
//     "home/index": {"layout/layout.html", "home/index.html"},
// }
func (r *Template) AddTemplate (fs fs.FS, filesMap map[string][]string) {
	for name, files := range filesMap {
		tName := path.Base(files[0])
		t, err := template.New(tName).Funcs(r.FuncMap).ParseFS(fs, files...)
		if err != nil {
			panic(err)
		}
		r.Render.Add(name, t)
	}
}


func NewTemplate() *Template {
	r := &Template{
		Render: multitemplate.NewRenderer(),
		FuncMap: make(map[string]interface{}),
	}

	// default funcMap
	r.AddFunc("unescaped", unescaped)

	return r
}

func unescaped(x string) interface{} {
	return template.HTML(x)
}
