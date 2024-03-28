package client

import (
	"github.com/jlaffaye/ftp"
	"io"
	"os"
	"strings"
)

type Client struct {
	cl      *ftp.ServerConn
	saveDir string
}

func NewClient(clientFTP *ftp.ServerConn, saveDir string) (*Client, error) {
	err := os.MkdirAll(saveDir, 0777)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(saveDir, "/") {
		saveDir += "/"
	}
	return &Client{
		cl:      clientFTP,
		saveDir: saveDir,
	}, nil
}

func (c *Client) ListFiles(path string) (string, error) {
	var s strings.Builder
	files, err := c.cl.List(path)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		s.WriteString(f.Name + "\n")
	}
	return s.String(), nil
}

func (c *Client) LoadFile(fileName string) error {
	f, err := os.Create(c.saveDir + fileName)
	if err != nil {
		return err
	}
	r, err := c.cl.Retr(fileName)
	if err != nil {
		return err
	}
	defer r.Close()
	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendFile(fileName, savePath string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	err = c.cl.Stor(savePath, f)
	return err
}

func (c *Client) Mkdir(dirName string) error {
	err := c.cl.MakeDir(dirName)
	return err
}

func (c *Client) DeleteFile(fileName string) error {
	err := c.cl.Delete(fileName)
	return err
}
