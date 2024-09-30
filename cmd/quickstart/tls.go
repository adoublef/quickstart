package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	// "golang.org/x/crypto/acme/autocert"
)

func setupTLS(insecureSkip bool) (*tls.Config, error) {
	// https://www.statuscake.com/blog/serving-multiple-ssl-certificates-in-your-go-tests/
	// https://shaneutt.com/blog/golang-ca-and-signed-cert-go/
	// https://goenning.net/blog/free-and-automated-ssl-certificates-with-go/
	// https://github.com/caddyserver/certmagic
	// get the system cert pool
	root, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to return system cert pool: %w", err)
	}
	// todo: append local/provate certs
	// https://blog.cloudflare.com/exposing-go-on-the-internet/
	c := &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
		InsecureSkipVerify: insecureSkip,
		RootCAs:            root,
	}
	return c, nil
}
