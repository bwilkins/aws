package v4

import (
  "strings"
  "sort"
)

type CanonicalHeaders map[string]string

func (headers *CanonicalHeaders) String() string {
  canonicalised := make([]string, 0)
  for headerName, headerValue := range *headers {
    normalisedName := strings.TrimSpace(strings.ToLower(headerName))
    normalisedValue := strings.TrimSpace(headerValue)
    canonicalised = append(canonicalised, normalisedName + ":" + normalisedValue)
  }
  sort.Strings(canonicalised)

  return strings.Join(canonicalised, "\n") + "\n"
}
