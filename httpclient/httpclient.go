// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Cloud Security Client Go contributors
//
// SPDX-License-Identifier: Apache-2.0
package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sap/cloud-security-client-go/env"
)

// DefaultTLSConfig creates default tls.Config. Initializes SystemCertPool with cert/key from identity config.
//
// identity provides certificate and key
func DefaultTLSConfig(identity env.Identity) (*tls.Config, error) {
	if !identity.IsCertificateBased() {
		return nil, nil
	}
	certPEMBlock := []byte(identity.GetCertificate())
	keyPEMBlock := []byte(identity.GetKey())

	tlsCert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, fmt.Errorf("error creating x509 key pair for DefaultTLSConfig: %w", err)
	}
	tlsCertPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("error setting up cert pool for DefaultTLSConfig: %w", err)
	}
	ok := tlsCertPool.AppendCertsFromPEM(certPEMBlock)
	if !ok {
		return nil, errors.New("error adding certs to pool for DefaultTLSConfig")
	}
	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		RootCAs:      tlsCertPool,
		Certificates: []tls.Certificate{tlsCert},
	}
	return tlsConfig, nil
}

// DefaultHTTPClient
//
// tlsConfig required in case of cert-based identity config
func DefaultHTTPClient(tlsConfig *tls.Config) *http.Client {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	if tlsConfig != nil {
		client.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
			MaxIdleConns:    50,
		}
	}
	return client
}
