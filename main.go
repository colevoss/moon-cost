package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

type CurlFile struct {
	Method   string               `json:"method"`
	URL      string               `json:"url"`
	Requests map[string]ReqConfig `json:"requests"`
}

type Env map[string]string

type Config struct {
	Env    Env
	Params map[string]string
}

type ReqConfig struct {
	Params map[string]string `json:"params"`
	Body   *ReqBodyConfig    `json:"body"`
	Expect *ReqExpectConfig  `json:"expect"`
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

func main() {
	file, err := os.Open("./curl.json")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	dec := json.NewDecoder(file)

	var cfg CurlFile

	if err := dec.Decode(&cfg); err != nil {
		panic(cfg)
	}

	for k, req := range cfg.Requests {
		if len(req.Body.JSON) > 0 {
			fmt.Printf("req.Body: %v\n", string(req.Body.JSON))
		}

		data := Config{
			Env: map[string]string{
				"base": "hello.com",
			},
			Params: req.Params,
		}

		urlBuffer := bytes.Buffer{}
		urlTempl := template.Must(template.New("url").Parse(cfg.URL))

		if err := urlTempl.Execute(&urlBuffer, data); err != nil {
			panic(err)
		}

		fmt.Printf("%s url: %v\n", k, urlBuffer.String())
	}
}
