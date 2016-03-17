// Package rtorrent implements a client for rTorrent.
package rtorrent

import (
	"net/http"

	"github.com/kolo/xmlrpc"
)

// A Client is an rTorrent client.  It can be used to retrieve a
// variety of statistics from rTorrent.
type Client struct {
	Downloads *DownloadService

	xrc *xmlrpc.Client
}

// New creates a new Client using the input XML-RPC address and an optional
// transport.  If transport is nil, a default one will be used.
func New(addr string, transport http.RoundTripper) (*Client, error) {
	xrc, err := xmlrpc.NewClient(addr, transport)
	if err != nil {
		return nil, err
	}

	c := &Client{
		xrc: xrc,
	}

	c.Downloads = &DownloadService{c: c}

	return c, nil
}

// Close frees a Client's resources.
func (c *Client) Close() error {
	return c.xrc.Close()
}

// DownloadTotal retrieves the total number of downloaded bytes since
// rTorrent startup.
func (c *Client) DownloadTotal() (int, error) {
	return c.getInt("get_down_total", "")
}

// UploadTotal retrieves the total number of uploaded bytes since
// rTorrent startup.
func (c *Client) UploadTotal() (int, error) {
	return c.getInt("get_up_total", "")
}

// DownloadRate retrieves the current download rate in bytes from rTorrent.
func (c *Client) DownloadRate() (int, error) {
	return c.getInt("get_down_rate", "")
}

// UploadRate retrieves the current upload rate in bytes from rTorrent.
func (c *Client) UploadRate() (int, error) {
	return c.getInt("get_up_rate", "")
}

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

// DownloadRate retrieves the current download rate in bytes for a specific
// download, by its info-hash.
func (s *DownloadService) DownloadRate(infoHash string) (int, error) {
	return s.c.getInt("d.get_down_rate", infoHash)
}

// UploadRate retrieves the current upload rate in bytes for a specific
// download, by its info-hash.
func (s *DownloadService) UploadRate(infoHash string) (int, error) {
	return s.c.getInt("d.get_up_rate", infoHash)
}

// getInt retrieves an integer value from the specified XML-RPC method.
func (c *Client) getInt(method string, arg string) (int, error) {
	var send interface{}
	if arg != "" {
		send = arg
	}

	var v int
	err := c.xrc.Call(method, send, &v)
	return v, err
}

// getStringSlice retrieves a slice of string values from the specified
// XML-RPC method.
func (c *Client) getStringSlice(method string, args ...string) ([]string, error) {
	send := []interface{}{""}
	for _, a := range args {
		send = append(send, a)
	}

	var v []string
	err := c.xrc.Call(method, send, &v)
	return v, err
}
