package groups

import (
	"slices"

	"github.com/davidallendj/partitions/internal/partitions"
)

type Group struct {
	Name   string
	Labels []string
}

func (g *Group) GetNodeIDs(pm *partitions.DefaultManager) []string {
	foundNodes := []string{}
	for _, label := range g.Labels {
		nodeID := pm.LookupMember(label)
		if nodeID != nil {
			// check and make sure we're not duplicating node IDs
			if !slices.Contains(foundNodes, *nodeID) {
				foundNodes = append(foundNodes, *nodeID)
			}
		}
	}
	return foundNodes
}

func (g *Group) GetPartitions(pm *partitions.DefaultManager) []string {
	foundPartitions := []string{}
	for _, label := range g.Labels {
		partition := pm.LookupPartitionByMemberID(label)
		if partition != nil {
			if !slices.Contains(foundPartitions, partition.ID) {
				foundPartitions = append(foundPartitions, partition.ID)
			}
		}
	}
	return foundPartitions
}
