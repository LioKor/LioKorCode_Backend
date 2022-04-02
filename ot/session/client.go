package session

import "liokoredu/ot/selection"

type Client struct {
	Name      string              `json:"name"`
	Selection selection.Selection `json:"selection"`
}
