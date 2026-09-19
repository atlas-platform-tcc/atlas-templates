package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/dbz", dbzHandler)
	mux.HandleFunc("/", rootHandler)

	addr := ":" + getenv("PORT", "8080")
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// dbzHandler reports whether the platform-injected database is reachable. The Atlas
// platform wires the connection (DB_HOST/DB_PORT + POSTGRES_* from the Secret);
// running real queries is the application's job, so this only proves reachability.
func dbzHandler(w http.ResponseWriter, r *http.Request) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("no database configured\n"))
		return
	}
	addr := net.JoinHostPort(host, getenv("DB_PORT", "5432"))
	if err := dial(addr); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "database %s unreachable: %v\n", addr, err)
		return
	}
	fmt.Fprintf(w, "database %s reachable (db=%s user=%s)\n", addr,
		os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_USER"))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if host := os.Getenv("DB_HOST"); host != "" {
		fmt.Fprintf(w, "hello from atlas walking skeleton (database: %s)\n", host)
		return
	}
	_, _ = w.Write([]byte("hello from atlas walking skeleton\n"))
}

func dial(addr string) error {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return err
	}
	return conn.Close()
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
