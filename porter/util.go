package porter

import (
	"archive/zip"
	"bytes"
	"os"
	"strings"
)

func openZip(path string) (*zip.Reader, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
}

func createZip(path string) (*zip.Writer, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return zip.NewWriter(file), nil
}

func fileNameFromPath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.TrimSuffix(path, "/")
	path = strings.TrimSuffix(path, ".zip")

	return path[strings.LastIndex(path, "/")+1:]
}

func directoryPathFromPath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.TrimSuffix(path, "/")
	path = strings.TrimSuffix(path, fileNameFromPath(path))

	if !strings.HasSuffix(path, "/") {
		return path
	}
	return path[:len(path)-1]
}

// ignoreNone ignores no textures.
func ignoreNone(_ string) bool {
	return false
}
