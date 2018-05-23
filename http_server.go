package httpx

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/smartystreets/logging"
)

type HTTPServer struct {
	logger         *logging.Logger
	certificatePEM string
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
			MaxHeaderBytes: 1024 * 2,
		},
	}
}
func (this *HTTPServer) WithTLS(certificatePEM string, tlsConfig *tls.Config) *HTTPServer {
	if this == nil {
		return nil
	}

	if tlsConfig == nil {
		tlsConfig = &tls.Config{
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

	if strings.Contains(certificatePEM, "----BEGIN") {
		if cert, err := tls.X509KeyPair([]byte(certificatePEM), []byte(certificatePEM)); err == nil {
			tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
		} else {
			this.logger.Fatal("[ERROR] Unable to parse TLS certificate data: ", err)
		}
	} else {
		this.certificatePEM = certificatePEM
	}
	this.inner.TLSConfig = tlsConfig
	return this
}

func (this *HTTPServer) Listen() {
	if this == nil {
		return
	}

	this.logger.Printf("[INFO] Listening for web traffic on %s.\n", this.inner.Addr)
	if err := this.listen(); err == nil {
		return
	} else if err == http.ErrServerClosed {
		this.logger.Fatal("[INFO] Server shut down gracefully.")
	} else {
		this.logger.Fatal("[ERROR] Unable to listen to web traffic: ", err)
	}
}
func (this *HTTPServer) listen() error {
	if len(this.certificatePEM) == 0 && (this.inner.TLSConfig == nil || len(this.inner.TLSConfig.Certificates) == 0) {
		return this.inner.ListenAndServe()
	}

	return this.inner.ListenAndServeTLS(this.certificatePEM, this.certificatePEM)
}

func (this *HTTPServer) Shutdown(timeout time.Duration) error {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return this.inner.Shutdown(ctx)
}

func (this *HTTPServer) Close() {
	this.Shutdown(DefaultShutdownTimeout)
}

var DefaultShutdownTimeout = time.Second
