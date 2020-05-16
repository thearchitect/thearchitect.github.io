package resources

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"

	"github.com/thearchitect/thearchitect.github.io/server/resources/resource"
)

//go:generate go-bindata -pkg resources -nocompress -nomemcopy -o ./bindata.gen.go -prefix resources/ ./resources/...

func dockerfile() string {
	return _bindataDockerfile
}

func zshrc() string {
	return _bindataZshrc
}

func DockerContext() func() io.Reader {
	buf := bytes.NewBuffer([]byte{})

	//gw, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	//if err != nil {
	//	panic(err)
	//}
	//defer func() {
	//	if err := gw.Close(); err != nil {
	//		panic(err)
	//	}
	//}()

	tw := tar.NewWriter(buf)
	defer func() {
		if err := tw.Close(); err != nil {
			panic(err)
		}
	}()

	addFile := func(name string, data []byte) {
		if err := tw.WriteHeader(&tar.Header{
			Name: name,
			Size: int64(len(data)),
			Mode: 0777,
		}); err != nil {
			panic(err)
		}

		if _, err := io.Copy(tw, bytes.NewReader(data)); err != nil {
			panic(err)
		}
	}

	addFile("Dockerfile", []byte(dockerfile()))
	addFile("zshrc", []byte(zshrc()))

	if exe, err := os.Executable(); err != nil {
		panic(err)
	} else if data, err := ioutil.ReadFile(exe); err != nil {
		panic(err)
	} else {
		addFile("thearchitect", data)
	}

	if err := tw.Flush(); err != nil {
		panic(err)
	}

	//if err := gw.Flush(); err != nil {
	//	panic(err)
	//}

	data := buf.Bytes()

	return func() io.Reader {
		return bytes.NewReader(data)
	}
}

func IndexHTML(embed bool) (index, webapp resource.Resource) {
	webapp = WebAppJS()

	indexText := renderTemplate(
		_bindataIndexhtml,
		map[string]interface{}{
			"webapp": webapp.HTMLTag(embed),
		},
	)

	min := minifier()

	indexText, err := min.String(string(resource.ContentTypeHTML), indexText)
	if err != nil {
		panic(err)
	}

	index = resource.NewIndexResource(fmt.Sprintln(
		"",
		strings.Repeat("\n", 128),
		"",
		indexText,
	))

	return index, webapp
}

func WebAppJS() resource.Resource {
	return resource.NewJSResource(_bindataWebappjs)
}

////////////////////////////////////////////////////////////////
//// Minifier
////
func minifier() *minify.M {
	m := minify.New()
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepWhitespace:      true,
	})
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	return m
}

////////////////////////////////////////////////////////////////
//// Templating
////
func renderTemplate(text string, v interface{}) string {
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer([]byte{})

	if err := tmpl.Execute(buf, v); err != nil {
		panic(err)
	}

	return buf.String()
}
