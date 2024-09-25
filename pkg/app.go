package pkg

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

const (
	APP_NAME    = "ropen"
	APP_VERSION = "1.0.0"
)

type App struct {
	Config

	ips []string
}

func NewApp(cfgPath string, port int) (*App, error) {
	if err := LoadCfg(cfgPath); err != nil {
		debug("failed to load config: %v", err)
	}

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
	addr := fmt.Sprintf("%s:%d", a.ips[0], a.Config.Port)

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if perm, err := hasReadPermission(path); !perm || err != nil {
		return fmt.Errorf("no read permission on %s", path)
	}

	if info.IsDir() {
		debug("serving directory %s on: http://%s/",
			path, addr)

		http.Handle("/", http.FileServer(http.Dir(path)))
	} else {
		debug("serving file %s on: http://%s/%s",
			path, addr, filepath.Base(path))

		http.HandleFunc("/"+filepath.Base(path), fileHandlerHelper(path))
	}

	if err := http.ListenAndServe(addr, nil); err != nil {
		return fmt.Errorf("failed to start HTTP server: %v", err)
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
