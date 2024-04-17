package catalog_tags

import (
	"path/filepath"
	"strings"
)

func getDirectories(path string) []string {
	catalogs := filepath.Dir(path)
	parts := strings.Split(catalogs, string(filepath.Separator))
	return parts
}
