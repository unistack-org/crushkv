package crush

import (
	"strings"

	"github.com/unistack-org/crushkv/crushmap"
)

func makeNode(m *crushmap.Map, bucket *crushmap.Bucket, parent *TestingNode) *TestingNode {
	var child *TestingNode
	node := new(TestingNode)
	node.ID = bucket.TypeName + ":" + bucket.Name
	node.Type = uint16(bucket.TypeID)
	node.Weight = bucket.Weight
	node.Alg = bucket.Alg
	//node.Children = make([]Node, len(bucket.Items))
	node.Parent = parent
	if parent != nil {
		node.ID = parent.ID + "->" + node.ID
		//		parent.Weight += node.Weight
	}

	for _, item := range bucket.Items {
		childBucket := m.GetBucketByName(item.Name)
		if childBucket != nil {
			child = makeNode(m, childBucket, node)
		} else {
			idx := strings.Index(item.Name, ".")
			child = &TestingNode{
				ID:     item.Name,
				Type:   m.GetTypeIDByName(item.Name[:idx]),
				Weight: item.Weight,
				Parent: node,
			}
		}
		child.ID = node.ID + "->" + child.ID
		if parent != nil {
			parent.Weight += child.Weight
		}

		switch child.Alg {
		case "straw2":
			child.Selector = NewStraw2Selector(child)
		}

		node.Children = append(node.Children, child)
	}

	if node.Weight == 0 {
		for _, child := range node.Children {
			node.Weight += child.GetWeight()
		}
	}

	switch bucket.Alg {
	case "straw2":
		node.Selector = NewStraw2Selector(node)
	}

	return node
}
