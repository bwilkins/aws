package aws

import (
  "net/http"
  "io"
  "io/ioutil"
  "strings"
  "bytes"
  "crypto/sha256"
  "encoding/hex"
  "sort"
  "os"
  "github.com/bwilkins/aws/util"
  "time"
)

type Request struct {
  mRequest *http.Request
  mSigningHeaders *SigningHeaders
  mCanonicalHeaders *CanonicalHeaders
  now time.Time
}

func NewRequest(method, urlStr string, body io.Reader) (*Request, error) {
  r, e := http.NewRequest(method, urlStr, body)
  request := Request{r, nil, nil, time.Now().UTC()}
  return &request, e
}

func (request *Request) signingHeaders() *SigningHeaders {
  if request.mSigningHeaders == nil {
    headerNames := make(SigningHeaders, 0)
    for headerName, _ := range request.mRequest.Header {
      if len(strings.TrimSpace(headerName)) > 0 {
        headerNames = append(headerNames, headerName)
      }
    }
    sort.Strings(headerNames)
    request.mSigningHeaders = &headerNames
  }
  return request.mSigningHeaders
}

func (request *Request) canonicalHeaders() *CanonicalHeaders {
  if request.mCanonicalHeaders == nil {
    headers := make(CanonicalHeaders, 0)
    for _, headerName := range *request.signingHeaders() {
      normalisedName  := strings.TrimSpace(headerName)
      normalisedValue := strings.TrimSpace(strings.Join(request.mRequest.Header[headerName], " "))
      if len(normalisedName) > 0 && len(normalisedValue) > 0 {
        headers[normalisedName] = normalisedValue
      }
    }

    request.mCanonicalHeaders = &headers
  }

  return request.mCanonicalHeaders
}

func (request *Request) HashedPayload() string {
  body, _ := ioutil.ReadAll(request.mRequest.Body)
  request.mRequest.Body = ioutil.NopCloser(bytes.NewReader(body))
  hashed := sha256.Sum256(body)
  return hex.EncodeToString(hashed[:])
}


func (request *Request) CanonicalString() string {
  return strings.Join([]string{
    request.mRequest.Method,
    request.mRequest.URL.Path,
    request.mRequest.URL.RawQuery,
    request.canonicalHeaders().String(),
    request.signingHeaders().String(),
    request.HashedPayload(),
  }, "\n")
}

func(request *Request) generateCanonicalRequestHash() string {
  return util.HashString(request.CanonicalString())
}

func (request *Request) CredentialScopeString() string {
  return strings.Join([]string{
    request.HashingDate(),
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
  return strings.Join(request.mRequest.Header["X-Amz-Date"], "")
}

func (request *Request) StringToSign() string {
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
  kDate := util.HMAC_SHA256([]byte(aws_secret), request.HashingDate())
  kRegion := util.HMAC_SHA256(kDate, OpsWorksRegion)
  kService := util.HMAC_SHA256(kRegion, OpsWorksServiceName)
  return util.HMAC_SHA256(kService, "aws4_request")
}

func (request *Request) HashingDate() string {
  return request.now.Format("20060102")
}

func (request *Request) Sign() {
  signature := hex.EncodeToString( util.HMAC_SHA256(request.SigningKey(), request.StringToSign()) )
  request.mRequest.Header.Set("Authorization", OpsWorksSignatureAlgorithm +
    " Credential=" + request.CredentialString() +
    ", SignedHeaders=" + request.signingHeaders().String() +
    ", Signature=" + signature,
  )
}
