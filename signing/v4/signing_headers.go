package v4

import "strings"

type SigningHeaders []string


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

