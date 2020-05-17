package nginx

type Config struct {
	DirectoryAvailable          string
	DirectoryEnabled            string
	DirectoryAvailableChangeOld string
	DirectoryAvailableChangeNew string
	DirectoryEnabledChangeOld   string
	DirectoryEnabledChangeNew   string
	EnableSymlinks              bool
	RegenerateResources         bool
	DirectoryAvailableChanged   bool
	DirectoryEnabledChanged     bool
}
