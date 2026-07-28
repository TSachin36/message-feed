package api

import (
	"hash/crc32"
	"sort"
)

var shardAddresses = []string{
	"localhost:50051",
	"localhost:50052",
	"localhost:50053",
}

var hashRing []uint32
var hashRingNodes = make(map[uint32]string)

func buildHashRing() {

	hashRing = nil
	hashRingNodes = make(map[uint32]string)

	for _, address := range shardAddresses {

		hash := crc32.ChecksumIEEE(
			[]byte(address),
		)

		hashRing = append(
			hashRing,
			hash,
		)

		hashRingNodes[hash] = address
	}

	sort.Slice(
		hashRing,
		func(i, j int) bool {
			return hashRing[i] < hashRing[j]
		},
	)
}

func shardForUser(userID string) string {

	hash := crc32.ChecksumIEEE(
		[]byte(userID),
	)

	index := sort.Search(
		len(hashRing),
		func(i int) bool {
			return hashRing[i] >= hash
		},
	)

	if index == len(hashRing) {
		index = 0
	}

	return hashRingNodes[hashRing[index]]
}
