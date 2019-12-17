package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

const (
	signatureHeader = "x-xero-signature"
)

// WebhookAuthorizationMiddleware will check if the webhook request has the correct
// signature
func WebhookAuthorizationMiddleware(webhookSigningKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			signatureVal := r.Header.Get(signatureHeader)
			mac := hmac.New(sha256.New, []byte(webhookSigningKey))
			buf, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body.Close()

			_, err = mac.Write(buf)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			expectedMAC := mac.Sum(nil)
			if base64.StdEncoding.EncodeToString(expectedMAC) != signatureVal {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
			next.ServeHTTP(w, r)
		})
	}
}
