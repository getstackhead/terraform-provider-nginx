package nginx

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

func getNginxDirectories(filename string, m Config) (string, string) {
	pathAvailable := path.Join(m.DirectoryAvailable, filename)
	if !m.EnableSymlinks {
		return pathAvailable, ""
	}
	pathEnabled := path.Join(m.DirectoryEnabled, filename)
	return pathAvailable, pathEnabled
}

func CreateOrUpdateServerBlock(filename string, content string, m Config, markers map[string]interface{}, markersSplit map[string]interface{}) (string, error) {
	fullPathAvailable, _ := getNginxDirectories(filename, m)
	content = ReplaceMarkers(content, ProcessMarkers(markers, markersSplit))

	if err := ioutil.WriteFile(fullPathAvailable, []byte(content), 0744); err != nil {
		return "", err
	}
	return fullPathAvailable, nil
}

/// ReplaceMarkers replaces all markers in given content
func ReplaceMarkers(content string, markers map[string]string) string {
	var re *regexp.Regexp
	for key, value := range markers {
		quotedKey := regexp.QuoteMeta(key)
		re, _ = regexp.Compile("{#\\s*" + quotedKey + "\\s*#}") // {#marker#}
		content = re.ReplaceAllString(content, value)
		re, _ = regexp.Compile("{~\\s*" + quotedKey + "\\s*~}") // {~marker~}
		content = re.ReplaceAllString(content, value)
		re, _ = regexp.Compile("{\\*\\s*" + quotedKey + "\\s*\\*}") // {*marker*}
		content = re.ReplaceAllString(content, value)
	}
	return content
}

/// ProcessMarkers resolves array values into single string replaces
func ProcessMarkers(markers map[string]interface{}, markers_split map[string]interface{}) map[string]string {
	markers_split_keys := make([]string, 0, len(markers_split))
	for k := range markers_split {
		markers_split_keys = append(markers_split_keys, k)
	}

	output := make(map[string]string)
	for key, value := range markers {
		stringValue := value.(string)
		splitChar := markers_split[key]
		if splitChar == nil || splitChar.(string) == "" {
			output[key] = stringValue
			continue
		}

		// Split value by character
		for i, slice := range strings.Split(stringValue, splitChar.(string)) {
			output[fmt.Sprintf("%s[%d]", key, i)] = slice
		}
	}
	return output
}

func RemoveNginxServerBlock(filename string, m Config) error {
	fullPathAvailable, fullPathEnabled := getNginxDirectories(filename, m)
	// Remove symlink if exists
	if m.EnableSymlinks && FileExists(fullPathEnabled) {
		if err := os.Remove(fullPathEnabled); err != nil {
			return err
		}
	}
	// Remove configuration
	if err := os.Remove(fullPathAvailable); err != nil {
		return err
	}
	return nil
}

func ReadFile(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

func DisableServerBlock(filename string, m Config) error {
	if !m.EnableSymlinks {
		return nil
	}
	_, fullPathEnabled := getNginxDirectories(filename, m)
	return os.Remove(fullPathEnabled)
}

func EnableServerBlock(filename string, m Config) error {
	if !m.EnableSymlinks {
		return nil
	}
	fullPathAvailable, fullPathEnabled := getNginxDirectories(filename, m)
	if err := os.Symlink(fullPathAvailable, fullPathEnabled); err != nil {
		return err
	}
	return nil
}

func FileExists(path string) bool {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}
