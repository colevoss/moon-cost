package curl

import (
	"bytes"
	"encoding/json"
	"io"
	"text/template"
)

type Curl struct {
	Method   string               `json:"method"`
	URL      string               `json:"url"`
	Requests map[string]ReqConfig `json:"requests"`
	Headers  map[string]string    `json:"headers"`
}

type Config struct {
	Env    Env
	Params map[string]string
}

type ReqConfig struct {
	Params  map[string]string `json:"params"`
	Headers map[string]string `json:"headers"`
	Body    *ReqBodyConfig    `json:"body"`
	Expect  *ReqExpectConfig  `json:"expect"`
}

type ReqBodyConfig struct {
	File string          `json:"file"`
	JSON json.RawMessage `json:"json"`
	// TODO: FORM??
	// TODO: Some how do some sort of "RAW" type as well
}

type ReqExpectConfig struct {
	Status int `json:"status"`
}

func (c *Curl) Read(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(c)
}

func (c *Curl) JSON(data []byte) error {
	return json.Unmarshal(data, c)
}

func (c *Curl) ParsedURL(config Config) (string, error) {
	var buf bytes.Buffer

	templ, err := template.New("url").Parse(c.URL)

	if err != nil {
		return "", err
	}

	if err := templ.Execute(&buf, config); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (c *Curl) Req(name string) (ReqConfig, bool) {
	config, ok := c.Requests[name]

	return config, ok
}
