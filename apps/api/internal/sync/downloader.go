package sync

import "pocketpanel/api/internal/models"

type JARDownloader interface {
	DownloadJAR(version, destPath string) error
	ServerType() models.ServerType
}
