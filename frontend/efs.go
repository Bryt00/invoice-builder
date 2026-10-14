package frontend

import (
	"embed"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DistFS embeds all compiled production files from the dist directory.
//
//go:embed all:dist
var distFS embed.FS

// FS returns an fs.FS rooted directly at the dist directory.
func FS() fs.FS {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return distFS
	}
	return sub
}

// SPAHandler implements http.Handler to serve the Vue Single Page Application (SPA).
// It prefers files on the host disk (allowing live development/updates) and
// seamlessly falls back to the embedded filesystem.
type SPAHandler struct {
	embeddedFS fs.FS
	diskDir    string
}

// NewHandler constructs an SPAHandler.
// diskDir is checked first; if a file does not exist on disk, embeddedFS is used.
func NewHandler(diskDir string) *SPAHandler {
	return &SPAHandler{
		embeddedFS: FS(),
		diskDir:    diskDir,
	}
}

// openFile attempts to find a non-directory file on disk, then in the embedded filesystem.
func (h *SPAHandler) openFile(name string) (io.ReadSeekCloser, fs.FileInfo, error) {
	cleanName := filepath.Clean(name)

	// 1. Check local disk if directory exists
	if h.diskDir != "" {
		diskPath := filepath.Join(h.diskDir, cleanName)
		info, err := os.Stat(diskPath)
		if err == nil && !info.IsDir() {
			file, err := os.Open(diskPath)
			if err == nil {
				return file, info, nil
			}
		}
	}

	// 2. Check embedded filesystem
	if h.embeddedFS != nil {
		fsPath := filepath.ToSlash(cleanName)
		file, err := h.embeddedFS.Open(fsPath)
		if err == nil {
			info, err := file.Stat()
			if err == nil && !info.IsDir() {
				if rsc, ok := file.(io.ReadSeekCloser); ok {
					return rsc, info, nil
				}
			}
			file.Close()
		}
	}

	return nil, nil, os.ErrNotExist
}

// ServeHTTP serves static assets or falls back to index.html for Vue client-side routes.
func (h *SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Clean and normalize incoming path
	reqPath := strings.TrimPrefix(filepath.Clean(r.URL.Path), "/")
	if reqPath == "." || reqPath == "" {
		reqPath = "index.html"
	}

	// Try to open requested file
	file, info, err := h.openFile(reqPath)
	if err != nil {
		// File does not exist -> fall back to index.html for client-side routing
		file, info, err = h.openFile("index.html")
		if err != nil {
			http.Error(w, "Vue SPA frontend build not found. Please build the frontend first.", http.StatusNotFound)
			return
		}
		reqPath = "index.html"
	}
	defer file.Close()

	// Apply headers based on asset type
	switch {
	case reqPath == "index.html" || reqPath == "sw.js" || strings.HasSuffix(reqPath, "manifest.json") || strings.HasSuffix(reqPath, "manifest.webmanifest"):
		// Never cache HTML, service worker, or manifest so updates are recognized immediately
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	case strings.HasPrefix(reqPath, "assets/"):
		// Hashed Vite bundle chunks
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	default:
		// Static icons, images, and fonts
		w.Header().Set("Cache-Control", "public, max-age=86400")
	}

	// Explicit MIME types for PWA webmanifest and fonts
	ext := filepath.Ext(reqPath)
	if ext == ".webmanifest" {
		w.Header().Set("Content-Type", "application/manifest+json")
	} else if cType := mime.TypeByExtension(ext); cType != "" {
		w.Header().Set("Content-Type", cType)
	}

	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}
