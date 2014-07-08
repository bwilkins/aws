package aws

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
)

const (
  OpsWorksSignatureAlgorithm = "AWS4-HMAC-SHA256"
  OpsWorksHost = "opsworks.us-east-1.amazonaws.com"
  OpsWorksRegion = "us-east-1"
  OpsWorksServiceName = "opsworks"
  OpsWorksVersion = "2013-02-18"
  OpsworksTargetPrefix = "OpsWorks_20130218"
)

func Do(request *Request, v interface{}) error {
  request.mRequest.Header.Set("X-Amz-Target", OpsworksTargetPrefix + "." + "DescribeInstances")
  request.mRequest.Header.Set("Host", OpsWorksHost)
  request.mRequest.Header.Set("Content-Type", "application/x-amz-json-1.1")
  request.mRequest.Header.Set("X-Amz-Date", request.now.Format("20060102T150405Z"))
  request.Sign()
  client := http.DefaultClient

  response, err := client.Do(request.mRequest)

  if err != nil {
    fmt.Println("FAIL")
    fmt.Println(err.Error())
    return err
  }

  return unmarshal(request, response, v)
}


type ErrorJSON struct {
  Type string `json:"__type"`
  Message string `json:"message"`
}

func unmarshal(req *Request, res *http.Response, v interface{}) error {
  if res.StatusCode == http.StatusOK {
    fmt.Println("SUCCESS")
    b, err := ioutil.ReadAll(res.Body)
    if err != nil {
      return err
    }
    fmt.Printf("%s", b)
    return json.Unmarshal(b, v)
  }

  badBody, err := ioutil.ReadAll(res.Body)
  var badBodyParsed ErrorJSON
  if badBody != nil {
    //Assume badBody is JSON
    json.Unmarshal(badBody, &badBodyParsed)
  }
  if err != nil {
    return err
  }

  fmt.Printf("An error occurred! Code: %d\n", res.StatusCode)
  fmt.Printf("Body: %s\n", badBody)
  fmt.Printf("Body: %s\n", badBodyParsed.Message)


  //fmt.Printf("The Canonical String for this request was: '%s'\n", req.CanonicalString())
  fmt.Printf("The String-to-Sign for this request was: '%s'\n", req.StringToSign())
  fmt.Println(req.mRequest.Header)
  return nil
}

