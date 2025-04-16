package curl

import (
	"encoding/json"
	"io"
)

type Curl struct {
	Method   string             `json:"method"`
	URL      string             `json:"url"`
	Requests map[string]Request `json:"requests"`
	Headers  map[string]string  `json:"headers"`
	Query    map[string]string  `json:"query"`
}

type Params struct {
	Env    Env
	Params ReqestParams
}

type Request struct {
	Params  ReqestParams      `json:"params"`
	Headers map[string]string `json:"headers"`
	Body    Body              `json:"body"`
	Expect  *Expect           `json:"expect"`
	Query   map[string]string `json:"query"`
}

type ReqestParams map[string]string

type Body struct {
	File string          `json:"file"`
	JSON json.RawMessage `json:"json"`
	// TODO: FORM??
	// TODO: Some how do some sort of "RAW" type as well
}

type Expect struct {
	Status int `json:"status"`
}

func (c *Curl) Read(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(c)
}

func (c *Curl) Req(name string) (Request, bool) {
	config, ok := c.Requests[name]

	return config, ok
}
