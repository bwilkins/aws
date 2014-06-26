package aws

import (
  "bytes"
  "crypto/hmac"
  "crypto/sha256"
  "encoding/hex"
  "encoding/json"
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  //"net/url"
  "os"
  "sort"
  "strings"
  "time"
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
  request.Header.Set("X-Amz-Target", OpsworksTargetPrefix + "." + "DescribeInstances")
  request.Header.Set("Host", OpsWorksHost)
  request.Header.Set("Content-Type", "application/x-amz-json-1.1")
  request.Header.Set("X-Amz-Date", time.Now().UTC().Format("20060102T150405Z"))
  request.Sign()
  client := http.DefaultClient

  response, err := client.Do((*http.Request)(request))

  if err != nil {
    fmt.Println("FAIL")
    fmt.Println(err.Error())
    return err
  }

  return unmarshal(request, response, v)
}

type Request http.Request

func NewRequest(method, urlStr string, body io.Reader) (*Request, error) {
  r, e := http.NewRequest(method, urlStr, body)
  request := (*Request)(r)
  return request, e
}

type SigningHeaders []string

func (request *Request) signingHeaders() *SigningHeaders {
  headerNames := make(SigningHeaders, 0)
  for headerName, _ := range request.Header {
    if len(strings.TrimSpace(headerName)) > 0 {
      headerNames = append(headerNames, headerName)
    }
  }

  sort.Strings(headerNames)

  return &headerNames
}

func (headers *SigningHeaders) String() string {
  hdrs := make([]string, 0)
  for _, headerName := range *headers {
    hdrName := strings.TrimSpace(strings.ToLower(headerName))
    if len(hdrName) > 0 {
      hdrs = append(hdrs, hdrName)
    }
  }
  return strings.Join(hdrs, ";")
}

type CanonicalHeaders map[string]string

func (request *Request) canonicalHeaders() *CanonicalHeaders {
  headers := make(CanonicalHeaders, 0)
  for _, headerName := range *request.signingHeaders() {
    normalisedName  := strings.TrimSpace(headerName)
    normalisedValue := strings.TrimSpace(strings.Join(request.Header[headerName], " "))
    if len(normalisedName) > 0 && len(normalisedValue) > 0 {
      headers[normalisedName] = normalisedValue
    }
  }

  return &headers
}

func (headers *CanonicalHeaders) String() string {
  canonicalised := make([]string, 0)
  for headerName, headerValue := range *headers {
    normalisedName := strings.TrimSpace(strings.ToLower(headerName))
    normalisedValue := strings.TrimSpace(headerValue)
    canonicalised = append(canonicalised, normalisedName + ":" + normalisedValue)
  }

  return strings.Join(canonicalised, "\n") + "\n"
}


func (request *Request) HashedPayload() string {
  body, _ := ioutil.ReadAll(request.Body)
  request.Body = ioutil.NopCloser(bytes.NewReader(body))
  hashed := sha256.Sum256(body)
  return hex.EncodeToString(hashed[:])
}

func hashString(to_hash string) string {
  hashed := sha256.Sum256([]byte(to_hash))
  return hex.EncodeToString(hashed[:])
}

func (request *Request) generateCanonicalRequest() string {
  return strings.Join([]string{
    request.Method,
    request.URL.Path,
    request.URL.RawQuery,
    request.canonicalHeaders().String(),
    request.signingHeaders().String(),
    request.HashedPayload(),
  }, "\n")
}

func(request *Request) generateCanonicalRequestHash() string {
  return hashString(request.generateCanonicalRequest())
}

func (request *Request) CredentialScopeString() string {
  return strings.Join([]string{
    time.Now().UTC().Format("20060102"),
    OpsWorksRegion,
    OpsWorksServiceName,
    "aws4_request",
  }, "/")
}

func (request *Request) CredentialString() string {
  return strings.Join([]string{
    os.Getenv("AWS_ACCESS_KEY_ID"),
    request.CredentialScopeString(),
  }, "/")
}

func (request *Request) AmazonDateString() string {
  return strings.TrimSpace(strings.Join(request.Header["X-Amz-Date"], ""))
}

func (request *Request) SigningString() string {
  return strings.Join([]string{
    OpsWorksSignatureAlgorithm,
    request.AmazonDateString(),
    request.CredentialScopeString(),
    request.generateCanonicalRequestHash(),
  }, "\n")
}

func (request *Request) SigningKey() []byte {
  secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
  aws_secret := "AWS4" + secret
  kDate := HMAC_SHA256([]byte(aws_secret), request.AmazonDateString())
  kRegion := HMAC_SHA256(kDate, OpsWorksRegion)
  kService := HMAC_SHA256(kRegion, OpsWorksServiceName)
  return HMAC_SHA256(kService, "aws4_request")
}

func HMAC_SHA256(key []byte, data string) []byte {
  h := hmac.New(sha256.New, key)
  h.Write([]byte(data))
  return h.Sum([]byte{})
}

func (request *Request) Sign() {
  signature := hex.EncodeToString( HMAC_SHA256(request.SigningKey(), request.SigningString()) )
  request.Header.Set("Authorization", OpsWorksSignatureAlgorithm +
  " Credential=" + request.CredentialString() +
  ", SignedHeaders=" + request.signingHeaders().String() +
  ", Signature=" + signature,
)
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
  _, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return err
  }

  fmt.Printf("An error occurred! Code: %d\n", res.StatusCode)
  fmt.Println(req.Header)
  return nil
}

// Utils

// Used for debugging
type logReader struct {
  r io.Reader
}

func (lr *logReader) Read(b []byte) (n int, err error) {
  n, err = lr.r.Read(b)
  fmt.Print(string(b))
  return
}

// Sugar
type DescribeInstancesResponse struct {
  Instances []Instance
}

type Instance struct {
  AmiId string
  Status  string
  Hostname    string
  IpAddress  string
}

func DescribeInstances() (*DescribeInstancesResponse, error) {
  bodyEncoded:=""

  r, _ := NewRequest("POST", "https://"+OpsWorksHost+"/", strings.NewReader(bodyEncoded))

  v := new(DescribeInstancesResponse)
  return v, Do(r, v)
}
