package curl

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

type ReqBodyType string

const (
	ReqBodyTypeFile ReqBodyType = "FILE"
	ReqBodyTypeJSON ReqBodyType = "JSON"
	ReqBodyTypeNone ReqBodyType = "NONE"
)

type ReqBody struct {
	JSON io.Reader
	File io.ReadCloser
	Type ReqBodyType
}

func NoneReqBody() ReqBody {
	return ReqBody{
		Type: ReqBodyTypeNone,
	}
}

func JSONReqBody(data json.RawMessage) (ReqBody, error) {
	var body ReqBody
	compacted := &bytes.Buffer{}

	if err := json.Compact(compacted, data); err != nil {
		return body, err
	}

	body.JSON = compacted
	body.Type = ReqBodyTypeJSON

	return body, nil
}

func FileReqBody(path string) (ReqBody, error) {
	var body ReqBody

	file, err := os.Open(path)

	if err != nil {
		return body, err
	}

	body.File = file
	body.Type = ReqBodyTypeFile

	return body, nil
}

func (r *ReqBody) Data() io.Reader {
	switch r.Type {
	case ReqBodyTypeNone:
		return nil
	case ReqBodyTypeFile:
		return r.File
	case ReqBodyTypeJSON:
		return r.JSON
	}

	return nil
}

func (r *ReqBody) Close() {
	if r.Type == ReqBodyTypeFile {
		r.File.Close()
	}
}
