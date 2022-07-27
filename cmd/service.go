package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	// envs can be found not only in file
	_ = godotenv.Load()
}

// name of environment variable
const (
	certFileEnv = "CERT_FILE"
	keyFileEnv  = "KEY_FILE"
	port        = "AUTH_SERVICE_PORT"
)

// envs gathers expected environment variables,
// returns error if any of env vars is not set.
func envs(envs ...string) (map[string]string, error) {
	em := make(map[string]string, len(envs))
	var ok bool
	for _, env := range envs {
		if em[env], ok = os.LookupEnv(env); !ok {
			log.Println(em[env])
			return nil, fmt.Errorf("environment variable %q must be set", env)
		}
	}
	return em, nil
}

// newTLSServer returns preconfigured TLS 1.3 server.
func newTLSServer(addr string, mux http.Handler) *http.Server {
	return &http.Server{
		Addr: addr,
		// Handler: api.Router(),
		Handler:           mux,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
			// avoid the use of P-384 and P-521. P-256 is immune
			// to timing attacks, whereas P-384 and P-521 are not.
			CurvePreferences: []tls.CurveID{tls.CurveP256},
		},
	}
}

func main() {
	em, err := envs(certFileEnv, keyFileEnv, port)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Proudly served with Go and HTTPS!\n"))
	})

	srv := newTLSServer(em[port], mux)

	log.Printf("server listening on: localhost%s...", srv.Addr)

	go func() {
		if err := srv.ListenAndServeTLS(em[certFileEnv], em[keyFileEnv]); err != http.ErrServerClosed {
			log.Fatal(err)
		} else {
			log.Println(err) // prints 'server closed'
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-stop
	log.Println("got os signal", s)

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}
