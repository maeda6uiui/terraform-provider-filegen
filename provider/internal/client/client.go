package client

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

type Client struct {
}

type CreateFilledFileResponse struct {
	Md5sum string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (c *Client) GetFileMd5sum(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	sum := hash.Sum(nil)
	return fmt.Sprintf("%x", sum), nil
}

func (c *Client) CalcFilledBytesMd5sum(size int, filler byte) string {
	hash := md5.New()

	chunk_size := min(size, 32*1024)
	chunk := make([]byte, chunk_size)
	for i := range chunk {
		chunk[i] = filler
	}

	for remaining := size; remaining > 0; remaining -= chunk_size {
		hash.Write(chunk[:min(remaining, chunk_size)])
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (c *Client) CreateFilledFile(filename string, size int, filler byte) (*CreateFilledFileResponse, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(file)
	for range size {
		_, err = writer.Write([]byte{filler})
		if err != nil {
			return nil, err
		}
	}
	if err := writer.Flush(); err != nil {
		return nil, err
	}
	file.Close()

	file_hash, err := c.GetFileMd5sum(filename)
	if err != nil {
		return nil, err
	}

	return &CreateFilledFileResponse{
		Md5sum: file_hash,
	}, nil
}

func (c *Client) DeleteFile(filename string) error {
	if err := os.Remove(filename); err != nil {
		return err
	}

	return nil
}
