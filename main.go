package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

func main() {
	resp, err := http.Get("https://pbs.twimg.com/profile_images/434324230486757376/_NJDCzqq_400x400.png")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	sourceData := bytes.NewBuffer(body)
	originalImage, err := png.Decode(sourceData)
	if err != nil {
		log.Fatal(err)
	}
	sizes := []uint{16, 32, 64, 128, 256}
	var sonotsDataUriSchemaList []string

	for _, size := range sizes {
		thumbnail := resize.Thumbnail(size, size, originalImage, resize.Bicubic)
		var thumbnailData bytes.Buffer
		jpeg.Encode(&thumbnailData, thumbnail, nil)
		sonotsDataUriSchema := fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString(thumbnailData.Bytes()))
		sonotsDataUriSchemaList = append(sonotsDataUriSchemaList, sonotsDataUriSchema)
	}

	var opts struct {
		Sonotss []string
	}
	opts.Sonotss = sonotsDataUriSchemaList

	tpl := template.Must(template.New("mytemplate").Parse(`<!DOCTYPE html>
<html>
	<head><title>Awesome Sonots</title></head>
	<body>
		<h1>Awesome Sonots</h1>
		{{range $_, $uriDataSchema := .Sonotss}}
		<p><img src="{{ $uriDataSchema }}"></p>{{end}}
	</body>
</html>`))
	tpl.Execute(os.Stdout, opts)
}
