// Torrent management
// All Torrent management API methods are under "torrents",
// e.g.: /api/v2/torrents/methodName.
package qbt_apiv2

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strconv"

	errwrp "github.com/pkg/errors"
)

// BasicTorrent holds a basic torrent object from qbittorrent
type BasicTorrent struct {
	AddedOn       int     `json:"added_on"`
	Category      string  `json:"category"`
	CompletionOn  int64   `json:"completion_on"`
	Dlspeed       int     `json:"dlspeed"`
	Eta           int     `json:"eta"`
	ForceStart    bool    `json:"force_start"`
	Hash          string  `json:"hash"`
	Name          string  `json:"name"`
	NumComplete   int     `json:"num_complete"`
	NumIncomplete int     `json:"num_incomplete"`
	NumLeechs     int     `json:"num_leechs"`
	NumSeeds      int     `json:"num_seeds"`
	Priority      int     `json:"priority"`
	Progress      float64 `json:"progress"`
	Ratio         float64 `json:"ratio"`
	SavePath      string  `json:"save_path"`
	SeqDl         bool    `json:"seq_dl"`
	Size          int     `json:"size"`
	State         string  `json:"state"`
	SuperSeeding  bool    `json:"super_seeding"`
	Upspeed       int     `json:"upspeed"`
}

// Torrent holds a torrent object from qbittorrent
// with more information than BasicTorrent
type Torrent struct {
	AdditionDate           int     `json:"addition_date"`
	Comment                string  `json:"comment"`
	CompletionDate         int     `json:"completion_date"`
	CreatedBy              string  `json:"created_by"`
	CreationDate           int     `json:"creation_date"`
	DlLimit                int     `json:"dl_limit"`
	DlSpeed                int     `json:"dl_speed"`
	DlSpeedAvg             int     `json:"dl_speed_avg"`
	Eta                    int     `json:"eta"`
	LastSeen               int     `json:"last_seen"`
	NbConnections          int     `json:"nb_connections"`
	NbConnectionsLimit     int     `json:"nb_connections_limit"`
	Peers                  int     `json:"peers"`
	PeersTotal             int     `json:"peers_total"`
	PieceSize              int     `json:"piece_size"`
	PiecesHave             int     `json:"pieces_have"`
	PiecesNum              int     `json:"pieces_num"`
	Reannounce             int     `json:"reannounce"`
	SavePath               string  `json:"save_path"`
	SeedingTime            int     `json:"seeding_time"`
	Seeds                  int     `json:"seeds"`
	SeedsTotal             int     `json:"seeds_total"`
	ShareRatio             float64 `json:"share_ratio"`
	TimeElapsed            int     `json:"time_elapsed"`
	TotalDownloaded        int     `json:"total_downloaded"`
	TotalDownloadedSession int     `json:"total_downloaded_session"`
	TotalSize              int     `json:"total_size"`
	TotalUploaded          int     `json:"total_uploaded"`
	TotalUploadedSession   int     `json:"total_uploaded_session"`
	TotalWasted            int     `json:"total_wasted"`
	UpLimit                int     `json:"up_limit"`
	UpSpeed                int     `json:"up_speed"`
	UpSpeedAvg             int     `json:"up_speed_avg"`
}

// Tracker holds a tracker object from qbittorrent
type Tracker struct {
	Msg      string `json:"msg"`
	NumPeers int    `json:"num_peers"`
	Status   string `json:"status"`
	URL      string `json:"url"`
}

// WebSeed holds a webseed object from qbittorrent
type WebSeed struct {
	URL string `json:"url"`
}

// TorrentFile holds a torrent file object from qbittorrent
type TorrentFile struct {
	IsSeed       bool    `json:"is_seed"`
	Name         string  `json:"name"`
	Priority     int     `json:"priority"`
	Progress     float64 `json:"progress"`
	Size         int     `json:"size"`
	PieceRange   []int   `json:"piece_range"`
	Availability float64 `json:"availability"`
}

func (c *Client) AddNewTorrent(opt Optional) error {
	resp, err := c.postMultipartData("torrents/add", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	if err = RespBodyOk(resp.Body, ErrAddTorrnetfailed); err != nil {
		return err
	}
	return nil
}

func (c *Client) AddNewTorrentViaUrl(url, path string, tags ...string) error {
	ap, err := filepath.Abs(path)
	if err != nil {
		return errwrp.Wrapf(err, "cannot conv abs_path %s", path)
	}
	opt := Optional{
		"urls":     url,
		"savepath": ap,
	}
	fmt.Println(len(tags))
	if len(tags) > 0 {
		var ts string
		for _, t := range tags {
			ts += t + ","
		}
		ts = ts[:len(ts)-1]
		opt["tags"] = ts
	}
	err = c.AddNewTorrent(opt)
	return err
}

func (c *Client) TorrentList(opt Optional) ([]BasicTorrent, error) {
	resp, err := c.postXwwwFormUrlencoded("torrents/info", opt)

	err = RespOk(resp, err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bt := new([]BasicTorrent)
	err = json.Unmarshal(b, bt)
	if err != nil {
		return nil, err
	}
	return *bt, nil
}

func (c *Client) GetTorrentProperties(hash string) (Torrent, error) {
	resp, err := c.postXwwwFormUrlencoded("torrents/properties", Optional{
		"hash": hash,
	})
	err = RespOk(resp, err)
	if err != nil {
		return Torrent{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Torrent{}, err
	}
	t := new(Torrent)
	err = json.Unmarshal(b, t)
	if err != nil {
		return Torrent{}, err
	}
	return *t, nil
}

func (c *Client) GetTorrentContents(hash string, indexes ...int) ([]TorrentFile, error) {
	opt := Optional{
		"hash": hash,
	}
	if len(indexes) > 0 {
		var idxes string
		for _, idx := range indexes {
			idxes += strconv.Itoa(idx) + "|"
		}
		idxes = idxes[:len(idxes)-1]
		opt["indexes"] = idxes
	}

	resp, err := c.postXwwwFormUrlencoded("torrents/files", opt)
	err = RespOk(resp, err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	tf := new([]TorrentFile)
	err = json.Unmarshal(b, tf)
	if err != nil {
		return nil, err
	}
	return *tf, nil
}
