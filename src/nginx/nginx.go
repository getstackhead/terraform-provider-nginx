package nginx

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

func getNginxDirectories(filename string, m Config) (string, string) {
	pathAvailable := path.Join(m.DirectoryAvailable, filename)
	if !m.EnableSymlinks {
		return pathAvailable, ""
	}
	pathEnabled := path.Join(m.DirectoryEnabled, filename)
	return pathAvailable, pathEnabled
}

func CreateOrUpdateServerBlock(filename string, content string, m Config, markers map[string]interface{}) (string, error) {
	fullPathAvailable, _ := getNginxDirectories(filename, m)

	// Replace markers in content
	var re *regexp.Regexp
	for key, value := range ProcessMarkers(markers) {
		quotedKey := regexp.QuoteMeta(key)
		re, _ = regexp.Compile("{#\\s*" + quotedKey + "\\s*#}") // {#marker#}
		content = re.ReplaceAllString(content, value)
		re, _ = regexp.Compile("{~\\s*" + quotedKey + "\\s*~}") // {~marker~}
		content = re.ReplaceAllString(content, value)
		re, _ = regexp.Compile("{\\*\\s*" + quotedKey + "\\s*\\*}") // {*marker*}
		content = re.ReplaceAllString(content, value)
	}

	if err := ioutil.WriteFile(fullPathAvailable, []byte(content), 0744); err != nil {
		return "", err
	}
	return fullPathAvailable, nil
}

/// ProcessMarkers resolves array values into single string replaces
func ProcessMarkers(markers map[string]interface{}) map[string]string {
	output := make(map[string]string)
	for key, value := range markers {
		switch value.(type) {
		case []string:
			// Resolve array references if value is array
			for i, slice := range value.([]string) {
				output[fmt.Sprintf("%s[%d]", key, i)] = slice
			}
			break
		default:
			// Pass string as is
			output[key] = value.(string)
			break
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
