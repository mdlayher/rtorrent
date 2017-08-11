package rtorrent

import "net"
import "net/url"
import "strings"

const (
	// downloadList is used in methods which retrieve a list of downloads.
	downloadList = "download_list"
)

// A DownloadService is a wrapper for Client methods which operate on downloads.
type DownloadService struct {
	c *Client
}

// All retrieves a list of all downloads from rTorrent.
func (s *DownloadService) All() ([]string, error) {
	return s.c.getStringSlice(downloadList)
}

// Started retrieves a list of started downloads from rTorrent.
func (s *DownloadService) Started() ([]string, error) {
	return s.c.getStringSlice(downloadList, "started")
}

// Stopped retrieves a list of stopped downloads from rTorrent.
func (s *DownloadService) Stopped() ([]string, error) {
	return s.c.getStringSlice(downloadList, "stopped")
}

// Complete retrieves a list of complete downloads from rTorrent.
func (s *DownloadService) Complete() ([]string, error) {
	return s.c.getStringSlice(downloadList, "complete")
}

// Incomplete retrieves a list of incomplete downloads from rTorrent.
func (s *DownloadService) Incomplete() ([]string, error) {
	return s.c.getStringSlice(downloadList, "incomplete")
}

// Hashing retrieves a list of hashing downloads from rTorrent.
func (s *DownloadService) Hashing() ([]string, error) {
	return s.c.getStringSlice(downloadList, "hashing")
}

// Seeding retrieves a list of seeding downloads from rTorrent.
func (s *DownloadService) Seeding() ([]string, error) {
	return s.c.getStringSlice(downloadList, "seeding")
}

// Leeching retrieves a list of leeching downloads from rTorrent.
func (s *DownloadService) Leeching() ([]string, error) {
	return s.c.getStringSlice(downloadList, "leeching")
}

// Active retrieves a list of active downloads from rTorrent.
func (s *DownloadService) Active() ([]string, error) {
	return s.c.getStringSlice(downloadList, "active")
}

// BaseFilename retrieves the base filename shown in the rTorrent UI for a
// specific download, by its info-hash.
func (s *DownloadService) BaseFilename(infoHash string) (string, error) {
	return s.c.getString("d.get_base_filename", infoHash)
}

// TrackerDomain retrieves the domain name of the first tracker for a specific
// download, by its info-hash.
func (s *DownloadService) TrackerDomain(infoHash string) (string, error) {
	u, err := s.c.getStringAtIndex("t.get_url", infoHash, 0)
	if err != nil {
		return u, err
	}
	parts, err := url.Parse(u)
	if err != nil {
		return u, err
	}
	host := parts.Host
	if strings.Contains(host, ":") {
		h, _, err := net.SplitHostPort(parts.Host)
		if err == nil {
			host = h
		}
	}
	return host, err
}

// DownloadRate retrieves the current download rate in bytes for a specific
// download, by its info-hash.
func (s *DownloadService) DownloadRate(infoHash string) (int, error) {
	return s.c.getInt("d.get_down_rate", infoHash)
}

// DownloadTotal retrieves the current download total in bytes for a specific
// download, by its info-hash.
func (s *DownloadService) DownloadTotal(infoHash string) (int, error) {
	return s.c.getInt("d.get_down_total", infoHash)
}

// UploadRate retrieves the current upload rate in bytes for a specific
// download, by its info-hash.
func (s *DownloadService) UploadRate(infoHash string) (int, error) {
	return s.c.getInt("d.get_up_rate", infoHash)
}

// UploadTotal retrieves the current upload total in bytes for a specific
// download, by its info-hash.
func (s *DownloadService) UploadTotal(infoHash string) (int, error) {
	return s.c.getInt("d.get_up_total", infoHash)
}
