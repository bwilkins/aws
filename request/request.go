package request

import (
  "bytes"
  "encoding/json"
  "encoding/xml"
  "errors"
  "io"
  "net/http"
  "time"

  "github.com/bwilkins/aws"
  "github.com/bwilkins/aws/signing/v4"
)

type Request struct {
HttpRequest *http.Request
Action string
EndpointDefinition aws.EndpointDefinition
Timestamp time.Time
}

func NewRequest(method, action string, endpoint aws.EndpointDefinition, unformatted_request interface {}) (*Request, error) {
	body, err := formatRequestBody(unformatted_request, endpoint.DataInterchangeFormat)
	if err != nil {
		return nil, err
	}
	urlStr := "https://" + endpoint.Host + "/"
	r, e := http.NewRequest(method, urlStr, body)
	request := Request{r, action, endpoint, time.Now().UTC()}
	return &request, e
}

func (request *Request) Do(v interface {}) error {
  if request.EndpointDefinition.SignatureVersion == 4 {
    var signingRequest = v4.NewSigningRequest(request)
    return signingRequest.Do(v)
  }

	return errors.New("Unsupported Signature Version")
}

func formatRequestBody(unformatted_request interface{}, format string) (io.Reader, error) {
	var bodyEncoded []byte
	var err error

	if format == "json" {
		bodyEncoded, err = json.Marshal(unformatted_request)
	} else if format == "xml" {
		bodyEncoded, err = xml.Marshal(unformatted_request)
	} else {
		return nil, errors.New("Unsupported data interchange format: " + format)
	}

	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(bodyEncoded)

	return body, err
}
