package openapi

import (
	"bytes"
	"embed"
	"encoding/json"
	"io/fs"
	"strings"
	"text/template"
)

type swaggerInitializerUrl struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

//go:embed swagger-ui-5.26.2/dist
var SwaggerUIFS embed.FS

const SwaggerUIFolder = "swagger-ui-5.26.2/dist"

var swaggerInitializerTpl = `
window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    // url: "https://petstore.swagger.io/v2/swagger.json",
    urls: {{.Urls}},
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};
`

// ScanApiFiles scan /*/*.json file
func ScanApiFiles(fsys fs.FS) ([]swaggerInitializerUrl, error) {
	var urls []swaggerInitializerUrl

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".json") || d.IsDir() {
			return nil
		}

		urls = append(urls, swaggerInitializerUrl{
			Url:  path,
			Name: path,
		})
		return nil
	})
	return urls, err
}

// ParseTemplate parse swagger initializer file
func ParseTemplate(urls []swaggerInitializerUrl) string {
	urlsJson, err := json.Marshal(urls)
	if err != nil {
		panic(err)
	}
	tpl := template.Must(template.New("swaggerInitializerTpl").Parse(swaggerInitializerTpl))
	b := &bytes.Buffer{}
	err = tpl.Execute(b, map[string]interface{}{"Urls": string(urlsJson)})
	if err != nil {
		panic(err)
	}
	return b.String()
}
