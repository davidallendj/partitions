package partitions

import "github.com/davidallendj/partitions/internal/nodes"

type partition[T comparable] struct {
	ID      string
	Members []T
}

type nodePartition = partition[nodes.Node]
type defaultPartition = partition[string]
