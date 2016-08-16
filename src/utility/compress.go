package utility

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io/ioutil"
)

func GZipCompressData(rsc []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(rsc)
	w.Close()
	return b.Bytes()
}

func UnGZipCompressData(rsc []byte) ([]byte, error) {
	var b = bytes.NewReader(rsc)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	ret, err1 := ioutil.ReadAll(r)
	r.Close()
	return ret, err1
}

func CompressData(rsc []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(rsc)
	w.Close()
	return b.Bytes()
}

func UnCompressData(rsc []byte) ([]byte, error) {
	var b = bytes.NewReader(rsc)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	ret, err1 := ioutil.ReadAll(r)
	r.Close()
	return ret, err1
}
