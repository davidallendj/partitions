package partitions

import (
	"errors"
	"fmt"
	"slices"

	"github.com/davidallendj/partitions/internal/nodes"
)

var (
	ErrExists   = errors.New("partition exists")
	ErrNotFound = errors.New("partition not found")
)

type Manager[T comparable] struct {
	partitions []partition[T]
}

func (m *Manager[T]) CreatePartition(id string, members []T) error {
	// todo: deduplicate node ids

	// check if the partition already exists and fail if it does
	if m.ContainsPartition(id) {
		return ErrExists
	}

	// check the manager to see if there a nodes being added to multiple partitions
	for _, member := range members {
		if m.ContainsMember(member) {
			return fmt.Errorf("%v: %v", nodes.ErrExists, member)
		}
	}

	// no other issues at this point, so add the partition with nodes
	m.partitions = append(m.partitions, partition[T]{ID: id, Members: members})
	return nil
}

func (m *Manager[T]) AddNodeToPartition(partitionID string, member T) error {
	// do lookup to find existing partition
	partition := m.LookupPartitionByID(partitionID)
	if partition != nil {
		// check if node already exists and fail if it does
		if m.ContainsMember(member) {
			return fmt.Errorf("%v: %v", nodes.ErrExists, member)
		}
	} else {
		// no partition found, so return error
		return ErrNotFound
	}
	// add member to partition and update manager
	partition.Members = append(partition.Members, member)
	return nil
}

func (m *Manager[T]) LookupPartitionByID(id string) *partition[T] {
	// try and get index of partition with ID
	index := slices.IndexFunc(m.partitions, func(p partition[T]) bool {
		return p.ID == id
	})
	// found a partition, so return it
	if index >= 0 {
		return &m.partitions[index]
	}
	return nil
}

func (m *Manager[T]) LookupPartitionByMemberID(member T) *partition[T] {
	partitionIndex, _ := m.lookupMember(member)
	if partitionIndex >= 0 {
		return &m.partitions[partitionIndex]
	}
	return nil
}

func (m *Manager[T]) LookupMember(member T) *T {
	var (
		partitionIndex int
		memberIndex    int
	)
	partitionIndex, memberIndex = m.lookupMember(member)
	return m.getNodeFromPartition(partitionIndex, memberIndex)
}

func (m *Manager[T]) ContainsPartition(id string) bool {
	return m.LookupPartitionByID(id) != nil
}

func (m *Manager[T]) ContainsMember(member T) bool {
	var (
		partitionIndex int
		memberIndex    int
	)

	partitionIndex, memberIndex = m.lookupMember(member)
	return partitionIndex >= 0 && memberIndex >= 0
}

func (m *Manager[T]) GetPartitions() []partition[T] {
	return m.partitions
}

func (m *Manager[T]) GetPartitionIDs() []string {
	partitionIDs := []string{}
	for _, partition := range m.partitions {
		partitionIDs = append(partitionIDs, partition.ID)
	}
	return partitionIDs
}

func (m *Manager[T]) GetPartitionMembers() []T {
	members := []T{}
	for _, partition := range m.partitions {
		members = append(members, partition.Members...)
	}
	return members
}

func (m *Manager[T]) getNodeFromPartition(partitionIndex int, memberIndex int) *T {
	if partitionIndex >= 0 && memberIndex >= 0 {
		return &m.partitions[partitionIndex].Members[memberIndex]
	}
	return nil
}

func (m *Manager[T]) lookupMember(member T) (int, int) {
	var (
		partitionIndex int
		memberIndex    int
	)

	// check all partitions for nodes
	for _, partition := range m.partitions {
		memberIndex = slices.IndexFunc(partition.Members, func(testMember T) bool {
			return member == testMember
		})
		// we found the node in the partition so return
		if memberIndex >= 0 {
			return partitionIndex, memberIndex
		}
		partitionIndex += 1
	}
	// return negative values to indicate the node was not found
	return -1, -1
}

type NodeManager = Manager[nodes.Node]
type DefaultManager = Manager[string]
