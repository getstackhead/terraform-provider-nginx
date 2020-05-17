package nginx

import (
	"io/ioutil"
	"os"
	"path"
)

func getNginxDirectories(filename string, m Config) (string, string) {
	pathAvailable := path.Join(m.DirectoryAvailable, filename)
	if !m.EnableSymlinks {
		return pathAvailable, ""
	}
	pathEnabled := path.Join(m.DirectoryEnabled, filename)
	return pathAvailable, pathEnabled
}

func CreateOrUpdateVhost(filename string, content string, m Config) (string, error) {
	fullPathAvailable, _ := getNginxDirectories(filename, m)
	if err := ioutil.WriteFile(fullPathAvailable, []byte(content), 0744); err != nil {
		return "", err
	}
	return fullPathAvailable, nil
}

// move an existing vhost to another location.
// have to move symlink (if exists) before moving the main file!
func MoveNginxVhost(oldFileName string, newFileName string, m Config) (string, error) {
	oldFileAvailable, oldFileEnabled := getNginxDirectories(oldFileName, m)
	newFileAvailable, newFileEnabled := getNginxDirectories(newFileName, m)

	// Move symlink
	if m.EnableSymlinks && FileExists(oldFileEnabled) {
		if err := os.Rename(oldFileEnabled, newFileEnabled); err != nil {
			return oldFileAvailable, err
		}
	}
	if err := os.Rename(oldFileAvailable, newFileAvailable); err != nil {
		if m.EnableSymlinks {
			// Roll back symlink file
			_ = os.Rename(newFileEnabled, oldFileEnabled)
		}
		return oldFileAvailable, err
	}
	return newFileAvailable, nil
}

func RemoveNginxVhost(filename string, m Config) error {
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

func DisableVhost(filename string, m Config) error {
	if !m.EnableSymlinks {
		return nil
	}
	_, fullPathEnabled := getNginxDirectories(filename, m)
	return os.Remove(fullPathEnabled)
}

func EnableVhost(filename string, m Config) error {
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
