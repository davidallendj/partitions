package nodes

import "errors"

var (
	ErrExists = errors.New("node exists")
)

type Node struct {
	ID string
}
