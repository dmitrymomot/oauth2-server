package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// GetFileTypeByURI returns the file type of the given URI
func GetFileTypeByURI(uri string) string {
	ext := filepath.Ext(uri)
	if ext == "" {
		parsedUri, err := url.Parse(uri)
		if err != nil {
			return ""
		}
		ext = parsedUri.Query().Get("ext")
		if ext != "" {
			ext = parsedUri.Query().Get("format")
		}
	}
	ext = strings.Trim(ext, ".")

	switch ext {
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "mp4":
		return "video/mp4"
	case "mov":
		return "video/quicktime"
	case "mp3":
		return "audio/mpeg"
	case "flac":
		return "audio/flac"
	case "wav":
		return "audio/wav"
	case "glb":
		return "model/gltf-binary"
	case "gltf":
		return "model/gltf+json"
	case "html":
		return "text/html"
	case "js":
		return "application/javascript"
	case "css":
		return "text/css"
	case "json":
		return "application/json"
	case "xml":
		return "application/xml"
	case "svg":
		return "image/svg+xml"
	case "ico":
		return "image/x-icon"
	case "zip":
		return "application/zip"
	case "pdf":
		return "application/pdf"
	case "txt":
		return "text/plain"
	case "md":
		return "text/markdown"
	case "csv":
		return "text/csv"
	}

	return ""
}

// GetFileByPath returns the file bytes from the given path.
// If the path is a URL, it will download the file and return the bytes.
// If the path is a local file, it will read the file and return the bytes.
func GetFileByPath(path string) ([]byte, error) {
	if strings.Contains(path, "://") {
		if !strings.HasPrefix(path, "http") {
			return nil, fmt.Errorf("this url schema is not supported: %s", path)
		}

		b, err := DownloadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to download file from url: %w", err)
		}

		return b, nil
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from local disk: %w", err)
	}

	return b, nil
}

// DownloadFile downloads the file from the given URL and returns the bytes.
func DownloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return b, nil
}
