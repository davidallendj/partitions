package main

import (
	"encoding/json"

	"github.com/davidallendj/partitions/internal/groups"
	"github.com/davidallendj/partitions/internal/partitions"
	"github.com/rs/zerolog/log"
)

func main() {
	var (
		pm = partitions.DefaultManager{}
		n1 = "nid0001"
		n2 = "nid0002"
		n3 = "nid0003"
		p1 = "test1"
		p2 = "test2"
	)

	// create new partitions with partition manager while testing adding
	// the same partition multiple times
	unwrapError(pm.CreatePartition(p1, []string{n3}))
	unwrapError(pm.CreatePartition(p2, nil))
	unwrapError(pm.CreatePartition(p2, nil))

	// try and put the same node in multiple partitions which should cause error
	unwrapError(pm.AddNodeToPartition(p1, n1))
	unwrapError(pm.AddNodeToPartition(p1, n1))
	unwrapError(pm.AddNodeToPartition(p2, n2))

	// try and put the same node in multiple groups
	var (
		g1 = groups.Group{
			Name:   "group1",
			Labels: []string{n1, n2, "hello"},
		}
		g2 = groups.Group{
			Name:   "group2",
			Labels: []string{n1, n3, "world"},
		}
	)

	g1NodeIDs := ToJSON(g1.GetNodeIDs(&pm))
	g1PartitionIDs := ToJSON(g1.GetPartitions(&pm))
	g2NodeIDs := ToJSON(g2.GetNodeIDs(&pm))
	g2PartitionIDs := ToJSON(g2.GetPartitions(&pm))

	log.Info().Any("manager.partitions", pm.GetPartitions()).Msg("partition manager")
	log.Info().
		Any("group", g1).
		Any("found node IDs in manager", g1NodeIDs).
		Any("partitions containing found nodes", g1PartitionIDs).
		Msg("group 1")
	log.Info().
		Any("group", g2).
		Any("found node IDs in manager", g2NodeIDs).
		Any("partitions containing found nodes", g2PartitionIDs).
		Msg("group 2")
}

func unwrapError(err error) {
	if err != nil {
		log.Error().Err(err).Msg("something went wrong...")
	}
}

func ToJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
