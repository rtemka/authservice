package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	testcertEnv = "TEST_CERT_FILE"
	testkeyEnv  = "TEST_KEY_FILE"
)

const testport = ":34443"

var testcert = os.Getenv(testcertEnv)
var testkey = os.Getenv(testkeyEnv)

func TestTLSServer(t *testing.T) {

	if testcert == "" || testkey == "" {
		t.Skipf("TLSServer not found test certificates envs %q, %q, skipping...",
			testcertEnv, testkeyEnv)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := newTLSServer(testport, mux)
	defer srv.Close()
	go srv.ListenAndServeTLS(testcert, testkey)

	t.Run("self_signed_certs", func(t *testing.T) {

		cert, err := os.ReadFile(testcert)
		if err != nil {
			t.Fatalf("TLSServer = error: %v", err)
		}

		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(cert); !ok {
			t.Fatalf("TLSServer = error: %v", err)
		}

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: certPool,
				},
			},
		}

		resp, err := client.Get("https://localhost" + testport)
		if err != nil {
			t.Fatalf("TLSServer = error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("TLSServer = got status %d, want %d", resp.StatusCode, http.StatusOK)
		}
	})

	if err := srv.Shutdown(context.Background()); err != nil {
		t.Errorf("TLSServer = error: %v", err)
	}
}

func TestClientTLSGoogle(t *testing.T) {
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: 30 * time.Second},
		"tcp",
		"www.google.com:443",
		&tls.Config{
			CurvePreferences: []tls.CurveID{tls.CurveP256},
			MinVersion:       tls.VersionTLS12,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	state := conn.ConnectionState()
	t.Logf("TLS 1.%d", state.Version-tls.VersionTLS10)
	t.Log(tls.CipherSuiteName(state.CipherSuite))
	t.Log(state.VerifiedChains[0][0].Issuer.Organization[0])
	_ = conn.Close()
}
