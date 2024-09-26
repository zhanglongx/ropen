package pkg

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

const (
	APP_NAME    = "ropen"
	APP_VERSION = "1.0.1"
)

type App struct {
	Config

	ips []string
}

func NewApp(cfgPath string, port int) (*App, error) {
	LoadCfg(cfgPath)

	// override options from users
	if port > 0 {
		Cfg.Port = port
	}

	ips, err := getIPs(Cfg.PreferIPs)
	if err != nil {
		return nil, err
	}

	return &App{
		Config: Cfg,
		ips:    ips,
	}, nil
}

func (a *App) Run(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if perm, err := hasReadPermission(path); !perm || err != nil {
		return fmt.Errorf("no read permission on %s", path)
	}

	addr := fmt.Sprintf("%s:%d", a.ips[0], a.Config.Port)

	var cert tls.Certificate
	protocol := "http"
	if a.Config.Issuer.CAPath != "" && a.Config.Issuer.KeyPath != "" {
		if ca, err := NewCerts(a.Config.Issuer.CAPath,
			a.Config.Issuer.KeyPath); err == nil {
			if cert, err = ca.GenerateWebsiteCerts(a.ips[0]); err == nil {
				protocol = "https"
			} else {
				debug("failed to generate website certificates: %v", err)
			}
		} else {
			debug("failed to load CA certificates: %v", err)
		}
	}

	if info.IsDir() {
		fmt.Printf("serving directory %s on: %s://%s/",
			path, protocol, addr)

		noCacheHandler := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				h.ServeHTTP(w, r)
			})
		}

		http.Handle("/", noCacheHandler(http.FileServer(http.Dir(path))))
	} else {
		fmt.Printf("serving file %s on: %s://%s/%s",
			path, protocol, addr, filepath.Base(path))

		http.HandleFunc("/"+filepath.Base(path), fileHandlerHelper(path))
	}

	if protocol == "https" {
		server := &http.Server{
			Addr:    addr,
			Handler: nil,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}
		if err := server.ListenAndServeTLS("", ""); err != nil {
			return fmt.Errorf("failed to start HTTPS server: %v", err)
		}
	} else {
		if err := http.ListenAndServe(addr, nil); err != nil {
			return fmt.Errorf("failed to start HTTP server: %v", err)
		}
	}

	return nil
}

func hasReadPermission(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsPermission(err) {
			return false, nil
		}
		return false, err
	}

	defer file.Close()

	return true, nil
}

func fileHandlerHelper(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := filepath.Base(path)

		if !filepath.HasPrefix(r.URL.Path, "/"+fileName) {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", "application/octet-stream")

		http.ServeFile(w, r, path)
	}
}
