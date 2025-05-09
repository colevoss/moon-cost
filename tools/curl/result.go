package curl

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Result struct {
	Name string `json:"request"` // name of the request

	Status   int `json:"status"`
	Expected int `json:"expected,omitempty"`

	URL        string `json:"string"`
	Method     string `json:"method"`
	DurationMS int64  `json:"durationMS"`

	Duration time.Duration `json:"-"`

	Response *http.Response `json:"-"`
	Request  Request        `json:"-"`
	bodyData []byte
}

func NewResult(reqName string, res *http.Response, req Request, duration time.Duration) (Result, error) {
	var result Result

	result.Response = res
	result.Request = req
	result.URL = res.Request.URL.String()
	result.Method = res.Request.Method
	result.Status = res.StatusCode
	result.Expected = req.Expect.Status
	result.Name = reqName
	result.DurationMS = duration.Milliseconds()
	result.Duration = duration

	bodyData, err := io.ReadAll(res.Body)

	if err != nil {
		return result, err
	}

	result.bodyData = bodyData

	return result, nil
}

// Custom marshaller that will include the body after its ben parsed according to the response content type
func (r Result) MarshalJSON() ([]byte, error) {
	// need an alias of the result type so it doesn't recusively MarshalJSON
	type resultAlias Result
	type body struct {
		*resultAlias
		Body any    `json:"body"`
		Err  string `json:"error,omitempty"`
	}

	data := body{
		// cast result to alias type
		resultAlias: (*resultAlias)(&r),
	}

	dataBody, err := r.BodyData()

	if err != nil {
		data.Err = err.Error()
		return json.Marshal(data)
	}

	data.Body = dataBody

	return json.Marshal(&data)
}

func (r Result) BodyData() (any, error) {
	if r.Response.Header.Get("Content-Type") != "application/json" {
		return string(r.bodyData), nil
	}

	var body map[string]any

	err := json.Unmarshal(r.bodyData, &body)

	return body, err
}

func (r Result) BodyString() string {
	return string(r.bodyData)
}

func (r Result) BodyJson(data any) error {
	return json.Unmarshal(r.bodyData, data)
}
