package resources

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"os"
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
