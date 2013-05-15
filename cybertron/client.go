package cybertron

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
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

func (c *Client) List(url string, max int) ([]Revision, error) {
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

func (c *Client) Create(url string, rev int, archiveType string, body io.Reader) error {
	return c.multipartUpload("POST", url, rev, body, nil)
}

func (c *Client) Update(url string, from, to int, archiveType string, body io.Reader) error {
	header := make(http.Header)
	header.Add("X-Rev", strconv.Itoa(to))
	return c.multipartUpload("PATCH", url, from, body, header)
}

func (c *Client) Delete(url string, rev int) error {
	r, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s?rev=%d", c.Url, url, rev), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
}

func (c *Client) multipartUpload(verb string, url string, rev int, body io.Reader, header http.Header) error {
	pr, pw := io.Pipe()
	mpw := multipart.NewWriter(pw)

	r, err := http.NewRequest(verb, fmt.Sprintf("%s%s?rev=%d", c.Url, url, rev), pr)
	if err != nil {
		return err
	}
	if header != nil {
		r.Header = header
	}
	r.Header.Set("Content-Type", mpw.FormDataContentType())

	go func() {
		filew, err := mpw.CreateFormFile("file", "file")
		if err != nil {
			return
		}
		io.Copy(filew, body)

		mpw.Close()
		pw.Close()
	}()

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
}
