package crush

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/unistack-org/crushkv/crushmap"
)

func TestCrushStraw2(t *testing.T) {
	tree, err := makeStraw2Tree(t, "ssd")
	if err != nil {
		t.Fatal(err)
	}

	/*
		fmt.Printf("root %s %f\n", tree.GetID(), tree.GetWeight())
		for _, n := range tree.GetChildren() {
			if n.GetChildren() != nil {
				for _, k := range n.GetChildren() {
					fmt.Printf("children node %s %f\n", k.GetID(), k.GetWeight())
				}
			}
			fmt.Printf("node %s %f\n", n.GetID(), n.GetWeight())
		}
	*/
	nodes1 := Select(tree, 15, 2, 11, nil)
	for _, node := range nodes1 {
		for _, n := range node.GetChildren() {
			fmt.Printf("[STRAW2] For key %d got node : %#+v\n", 15, n.GetID())
		}
	}

	//	nodes2 := Select(tree, 4564564564, 2, 11, nil)
	/*
		nodes3 := Select(tree, 8789342322, 3, NODE, nil)
	*/
	/*
		for i := 0; i < 100000000; i++ {
			for _, node := range nodes1 {
				n := node.GetChildren()[0]
				if n.GetID() != "root:ssd->disktype:rmosd1_ssd->osd.3" {
					t.Fatal(fmt.Sprintf("[STRAW] For key %d got node : %#+v\n", 15, n))
				}
			}
		}
	*/
	/*
		for _, node := range nodes2 {
			n := node.GetChildren()[0]
			log.Printf("[STRAW] For key %d got node : %#+v\n", 4564564564, n)
		}
		/*
			for _, node := range nodes3 {
				log.Printf("[STRAW] For key %d got node : %s", 8789342322, node.GetID())
			}
	*/
}

func makeStraw2Tree(t *testing.T, pool string) (*TestingNode, error) {
	buf, err := ioutil.ReadFile("crushmap/testdata/map.txt2")
	if err != nil {
		return nil, err
	}

	m := crushmap.NewMap()
	err = m.DecodeText(buf)
	if err != nil {
		return nil, err
	}

	return makeNode(m, m.GetBucketByName(pool), nil), nil
}
