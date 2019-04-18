package crush

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"sort"
)

type HashingSelector struct {
	tokenList tokenList
	tokenMap  map[uint32]Node
}

func NewHashingSelector(n Node) *HashingSelector {
	var h = new(HashingSelector)
	var nodes = n.GetChildren()
	var maxWeight float32
	for _, node := range nodes {
		maxWeight = Max64(maxWeight, node.GetWeight())
	}
	h.tokenMap = make(map[uint32]Node)
	for _, node := range nodes {
		count := 500 * node.GetWeight() / maxWeight
		var hash []byte
		for i := float32(0); i < count; i++ {
			var input []byte
			if len(hash) == 0 {
				input = []byte(node.GetID())
			} else {
				input = hash
			}
			hash = digestBytes(input)
			token := Btoi(hash)
			if _, ok := h.tokenMap[token]; !ok {
				h.tokenMap[token] = node
			}
		}
	}
	h.tokenList = make([]uint32, 0, len(h.tokenMap))
	for k := range h.tokenMap {
		h.tokenList = append(h.tokenList, k)
	}
	sort.Sort(h.tokenList)
	return h

}

type tokenList []uint32

func (t tokenList) Len() int {
	return len(t)
}
func (t tokenList) Less(i, j int) bool {
	return t[i] < t[j]
}

func (t tokenList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func Max64(a float32, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func digestInt64(input int64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, input)
	bytes := buf.Bytes()
	result := sha1.Sum(bytes)
	hash := make([]byte, 20)
	copy(hash[:], result[:20])
	return hash
}
func digestBytes(input []byte) []byte {
	result := sha1.Sum(input)
	hash := make([]byte, 20)
	copy(hash[:], result[:20])
	return hash
}

func digestString(input string) []byte {
	result := sha1.Sum([]byte(input))
	hash := make([]byte, 20)
	copy(hash[:], result[:20])
	return hash
}

func Btoi(b []byte) uint32 {
	var result uint32
	buf := bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &result)
	return result
}

func (s *HashingSelector) Select(input uint32, round uint32) Node {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, input)
	binary.Write(buf, binary.LittleEndian, round)
	bytes := buf.Bytes()
	hash := digestBytes(bytes)
	token := Btoi(hash)
	return s.tokenMap[s.findToken(token)]
}

func (s *HashingSelector) findToken(token uint32) uint32 {
	i := sort.Search(len(s.tokenList), func(i int) bool { return s.tokenList[i] >= token })
	if i == len(s.tokenList) {
		return s.tokenList[i-1]
	}
	return s.tokenList[i]
}
