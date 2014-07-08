package util

import (
  "crypto/hmac"
  "crypto/sha256"
  "encoding/hex"
)

func HMAC_SHA256(key []byte, data string) []byte {
  h := hmac.New(sha256.New, key)
  h.Write([]byte(data))
  return h.Sum([]byte{})
}

func HashString(to_hash string) string {
  hashed := sha256.Sum256([]byte(to_hash))
  return hex.EncodeToString(hashed[:])
}
