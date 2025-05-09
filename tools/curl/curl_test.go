package curl

import (
	"strings"
	"testing"
)

func TestCurlRead(t *testing.T) {
	curlString := `
{
  "method": "POST",
  "url": "URL"
}
  `

	var curl Curl

	if err := curl.Read(strings.NewReader(curlString)); err != nil {
		t.Errorf("curl.Read() = %s. want nil", err)
	}

	if curl.Method != "POST" {
		t.Errorf("curl.Method = %s. want %s", curl.Method, "POST")
	}

	if curl.URL != "URL" {
		t.Errorf("curl.URL = %s. want %s", curl.URL, "URL")
	}
}

func TestCurlRequest(t *testing.T) {
	curl := Curl{
		Requests: map[string]Request{
			"a": {
				Body: Body{
					File: "a",
				},
			},
		},
	}

	aReq, aOk := curl.Req("a")

	if !aOk {
		t.Errorf("curl.Req(a) = _, %t. want %t", aOk, true)
	}

	if aReq.Body.File != "a" {
		t.Errorf("curl.Req(a).Body.File = %s. want %s", aReq.Body.File, "a")
	}

	_, bOk := curl.Req("b")

	if bOk {
		t.Errorf("curl.Req(b) = _, %t. want %t", bOk, false)
	}
}
