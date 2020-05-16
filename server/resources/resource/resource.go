package resource

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// type ReaderFunc func() io.Reader

type ContentType string

const (
	ContentTypeHTML ContentType = "text/html"
	ContentTypeCSS              = "text/css"
	ContentTypeJS               = "application/javascript"
)

////////////////////////////////////////////////////////////////
//// Resources
////

func newResource(data string, ct ContentType, cacheForever bool) Resource {
	res := fileResource{
		hash: hashString(data),
		data: []byte(data),
		ct:   ct,
		headers: http.Header{},
	}

	res.headers.Set("Content-Type", string(ct))
	res.headers.Set("Content-Length", fmt.Sprint(len(res.data)))

	if cacheForever {
		res.headers.Set("Cache-Control", "max-age=31536000")
	} else {
		res.headers.Set("Cache-Control", "no-store, must-revalidate")
		res.headers.Set("Pragma", "no-cache")
		res.headers.Set("Expires", "0")
	}

	return res
}

func NewIndexResource(data string) Resource {
	return newResource(data, ContentTypeHTML, false)
}

func NewJSResource(data string) Resource {
	return newResource(data, ContentTypeJS, true)
}

////////////////////////////////////////////////////////////////
//// Resource
////

type Resource interface {
	servable

	Data() []byte
	Hash() string

	HTMLTag(embed bool) string
}

type servable interface {
	http.Handler

	Mount(mux *http.ServeMux)
}

var _ Resource = new(fileResource)
var _ servable = new(fileResource)

type fileResource struct {
	hash string
	data []byte

	ct ContentType

	headers http.Header
}

func (res fileResource) Data() []byte {
	return res.data
}

func (res fileResource) Hash() string {
	return res.hash
}

func (res fileResource) HTMLTag(embed bool) string {
	switch res.ct {
	case ContentTypeJS:
		if embed {
			return fmt.Sprintf(
				`<script type="%s">%s</script>`,
				res.ct,
				string(res.data),
			)
		} else {
			return fmt.Sprintf(
				`<script type="%s" src="%s"></script>`,
				res.ct,
				res.mountPoint(),
			)
		}
	default:
		panic(errors.New(fmt.Sprintf("unsupported %s", res.ct)))
	}
}

func (res fileResource) applyHeaders(headers http.Header) {
	for hdr, val := range res.headers {
		for _, val := range val {
			headers.Add(hdr, val)
		}
	}
}

func (res fileResource) ServeHTTP(w http.ResponseWriter, q *http.Request) {
	res.applyHeaders(w.Header())

	if n, err := io.Copy(w, bytes.NewReader(res.data)); err != nil {
		panic(err)
	} else if n != int64(len(res.data)) {
		panic(errors.New(fmt.Sprintf("not enough bytes written n(%d) != len(%d)", n, len(res.data))))
	}
}

func (res fileResource) mountPoint() string {
	return fmt.Sprintf("/%s", res.hash)
}

func (res fileResource) Mount(mux *http.ServeMux) {
	mux.HandleFunc(
		res.mountPoint(),
		res.ServeHTTP,
	)
}

////////////////////////////////////////////////////////////////
//// Hashing
////

func hashString(v string) string {
	sum := sha512.Sum512([]byte(v))
	return base64.URLEncoding.EncodeToString(sum[:])
}
