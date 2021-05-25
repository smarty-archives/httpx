package httpx

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"strings"
	"time"
)

type HTTPServer struct {
	certificatePEM string
	keyPEM         string
	inner          *http.Server
}

func NewHTTPServer(listenAddress string, handler http.Handler) *HTTPServer {
	if len(listenAddress) == 0 {
		log.Println("[INFO] No listen address provided. No HTTP server will be created.")
		return nil
	}

	return &HTTPServer{
		inner: &http.Server{
			Addr:           listenAddress,
			Handler:        handler,
			ReadTimeout:    time.Second * 15,
			WriteTimeout:   time.Second * 15,
			MaxHeaderBytes: 1024 * 64,
		},
	}
}

func (this *HTTPServer) WithTLSFiles(cert, key string, tlsConfig *tls.Config) *HTTPServer {
	if this == nil {
		return nil
	}

	tlsConfig = defaultTLSConfig(tlsConfig)
	this.certificatePEM = cert
	this.keyPEM = key
	return this
}
func (this *HTTPServer) WithTLS(certificatePEM string, tlsConfig *tls.Config) *HTTPServer {
	if this == nil {
		return nil
	}

	tlsConfig = defaultTLSConfig(tlsConfig)
	if strings.Contains(certificatePEM, "----BEGIN") {
		if cert, err := tls.X509KeyPair([]byte(certificatePEM), []byte(certificatePEM)); err == nil {
			tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
		} else {
			log.Fatal("[ERROR] Unable to parse TLS certificate data: ", err)
		}
	} else {
		this.certificatePEM = certificatePEM
	}
	this.inner.TLSConfig = tlsConfig
	return this
}
func defaultTLSConfig(config *tls.Config) *tls.Config {
	if config != nil {
		return config
	}

	return &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		SessionTicketsDisabled:   true,
		CipherSuites: []uint16{
			tls.TLS_FALLBACK_SCSV,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

func (this *HTTPServer) Listen() {
	if this == nil {
		return
	}

	log.Printf("[INFO] Listening for HTTP traffic on %s.\n", this.inner.Addr)
	if err := this.listen(); err == nil {
		return
	} else if err == http.ErrServerClosed {
		log.Println("[INFO] HTTP listener shut down gracefully.")
	} else {
		log.Fatal("[ERROR] Unable to listen to HTTP traffic: ", err)
	}
}
func (this *HTTPServer) listen() error {
	if len(this.certificatePEM) == 0 && (this.inner.TLSConfig == nil || len(this.inner.TLSConfig.Certificates) == 0) {
		return this.inner.ListenAndServe()
	}

	return this.inner.ListenAndServeTLS(this.certificatePEM, this.privateKeyPEM())
}
func (this *HTTPServer) privateKeyPEM() string {
	if len(this.keyPEM) == 0 {
		return this.certificatePEM
	} else {
		return this.keyPEM
	}
}

func (this *HTTPServer) Shutdown(timeout time.Duration) error {
	if this == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return this.inner.Shutdown(ctx)
}

func (this *HTTPServer) Close() {
	this.Shutdown(DefaultShutdownTimeout)
}

var DefaultShutdownTimeout = time.Second
