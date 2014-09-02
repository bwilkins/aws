package v4

import (
  "net/http"
  "io/ioutil"
  "strings"
  "bytes"
  "crypto/sha256"
  "encoding/hex"
  "sort"
  "github.com/bwilkins/aws"
  "github.com/bwilkins/aws/signing/util"
  "time"
  "encoding/json"
  "fmt"
)

const (
  SignatureAlgorithm = "AWS4-HMAC-SHA256"
  SignatureVersionString = "AWS4"
)
type Request struct {
  HttpRequest *http.Request
  Action string
  EndpointDefinition aws.EndpointDefinition
  Timestamp time.Time
}

type SigningRequest struct {
  HttpRequest *http.Request
  Action string
  EndpointDefinition aws.EndpointDefinition
  Timestamp time.Time
  mSigningHeaders *SigningHeaders
  mCanonicalHeaders *CanonicalHeaders
}

func NewSigningRequest(request interface{}) SigningRequest {
  cast_request, _ := (request).(Request)
  return SigningRequest{
	cast_request.HttpRequest,
	cast_request.Action,
	cast_request.EndpointDefinition,
	cast_request.Timestamp,
	nil,
	nil,
  }
}

func (request *SigningRequest) signingHeaders() *SigningHeaders {
  if request.mSigningHeaders == nil {
    headerNames := make(SigningHeaders, 0)
    for headerName, _ := range request.HttpRequest.Header {
      if len(strings.TrimSpace(headerName)) > 0 {
        headerNames = append(headerNames, headerName)
      }
    }
    sort.Strings(headerNames)
    request.mSigningHeaders = &headerNames
  }
  return request.mSigningHeaders
}

func (request *SigningRequest) canonicalHeaders() *CanonicalHeaders {
  if request.mCanonicalHeaders == nil {
    headers := make(CanonicalHeaders, 0)
    for _, headerName := range *request.signingHeaders() {
      normalisedName  := strings.TrimSpace(headerName)
      normalisedValue := strings.TrimSpace(strings.Join(request.HttpRequest.Header[headerName], " "))
      if len(normalisedName) > 0 && len(normalisedValue) > 0 {
        headers[normalisedName] = normalisedValue
      }
    }

    request.mCanonicalHeaders = &headers
  }

  return request.mCanonicalHeaders
}

func (request *SigningRequest) HashedPayload() string {
  body, _ := ioutil.ReadAll(request.HttpRequest.Body)
  request.HttpRequest.Body = ioutil.NopCloser(bytes.NewReader(body))
  hashed := sha256.Sum256(body)
  return hex.EncodeToString(hashed[:])
}


func (request *SigningRequest) CanonicalString() string {
  return strings.Join([]string{
    request.HttpRequest.Method,
    request.HttpRequest.URL.Path,
    request.HttpRequest.URL.RawQuery,
    request.canonicalHeaders().String(),
    request.signingHeaders().String(),
    request.HashedPayload(),
  }, "\n")
}

func(request *SigningRequest) generateCanonicalRequestHash() string {
  return util.HashString(request.CanonicalString())
}

func (request *SigningRequest) CredentialScopeString() string {
  return strings.Join([]string{
    request.HashingDate(),
    request.EndpointDefinition.Region,
    request.EndpointDefinition.ServiceName,
    "aws4_request",
  }, "/")
}

func (request *SigningRequest) CredentialString() string {
  return strings.Join([]string{
    aws.AccessCredentials.AccessId,
    request.CredentialScopeString(),
  }, "/")
}

func (request *SigningRequest) AmazonDateString() string {
  return strings.Join(request.HttpRequest.Header["X-Amz-Date"], "")
}

func (request *SigningRequest) StringToSign() string {
  return strings.Join([]string{
    SignatureAlgorithm,
    request.AmazonDateString(),
    request.CredentialScopeString(),
    request.generateCanonicalRequestHash(),
  }, "\n")
}

func (request *SigningRequest) SigningKey() []byte {
  secret := aws.AccessCredentials.SecretKey
  aws_secret := SignatureVersionString + secret
  kDate := util.HMAC_SHA256([]byte(aws_secret), request.HashingDate())
  kRegion := util.HMAC_SHA256(kDate, request.EndpointDefinition.Region)
  kService := util.HMAC_SHA256(kRegion, request.EndpointDefinition.ServiceName)
  return util.HMAC_SHA256(kService, "aws4_request")
}

func (request *SigningRequest) HashingDate() string {
  return request.Timestamp.Format("20060102")
}

func (request *SigningRequest) Sign() {
  signature := hex.EncodeToString( util.HMAC_SHA256(request.SigningKey(), request.StringToSign()) )
  request.HttpRequest.Header.Set("Authorization", SignatureAlgorithm +
    " Credential=" + request.CredentialString() +
    ", SignedHeaders=" + request.signingHeaders().String() +
    ", Signature=" + signature,
  )
}

func (request *SigningRequest) Do(v interface{}) error {
  request.HttpRequest.Header.Set("X-Amz-Target", request.EndpointDefinition.TargetPrefix + "." + request.Action)
  request.HttpRequest.Header.Set("Host", request.EndpointDefinition.Host)
  request.HttpRequest.Header.Set("Content-Type", "application/x-amz-json-1.1")
  request.HttpRequest.Header.Set("X-Amz-Date", request.Timestamp.Format("20060102T150405Z"))
  request.Sign()
  client := http.DefaultClient

  response, err := client.Do(request.HttpRequest)

  if err != nil {
    fmt.Println("FAIL")
    fmt.Println(err.Error())
    return err
  }

  return unmarshal(request, response, v)
}


func unmarshal(req *SigningRequest, res *http.Response, v interface{}) error {
  if res.StatusCode == http.StatusOK {
    b, err := ioutil.ReadAll(res.Body)
    if err != nil {
      return err
    }
    return json.Unmarshal(b, v)
  }

  _, bodyReadErr := ioutil.ReadAll(res.Body)
  if bodyReadErr != nil {
    return bodyReadErr
  }

  return nil
}
