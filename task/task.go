package mq2task

import (
	"io/ioutil"

	"github.com/pkg/errors"
)

// Client represents a task import
type Client struct {
}

// New creates a new Task import
func New(filename string) (c *Client, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		err = errors.Wrapf(err, "failed to read %s", filename)
		return
	}
	c = &Client{}
	return
}
