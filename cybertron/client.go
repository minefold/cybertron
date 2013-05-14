package cybertron

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

const (
	TarGz = "application/x-gzip"
)

type Revision struct {
	Rev int
	At  time.Time
}

type Client struct {
	Url string
}

func NewClient(url string) *Client {
	return &Client{Url: url}
}

func (c *Client) Create(url string, revision int, archiveType string, body io.Reader) error {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	fw, err := w.CreateFormFile("file", "localfile.tar.gz")
	if err != nil {
		return err
	}

	_, err = io.Copy(fw, body)
	if err != nil {
		return err
	}
	w.Close()

	req, err := http.NewRequest("POST", c.Url+url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	// req.SetBasicAuth("email@email.com", "password")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return err
}

func (c *Client) ListRevs(url string, max int) ([]Revision, error) {
	resp, err := http.Get(c.Url + url + ".json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	revs := make([]Revision, 0)
	if err := dec.Decode(&revs); err != nil {
		return nil, err
	}

	return revs, nil
}
