package crush

type Selector interface {
	Select(input uint32, round uint32) Node
}

type Node interface {
	GetChildren() []Node
	GetType() uint16
	GetWeight() float32
	GetID() string
	IsFailed() bool
	GetSelector() Selector
	SetSelector(Selector)
	GetParent() Node
	IsLeaf() bool
	Select(input uint32, round uint32) Node
}

type Comparitor func(Node) bool

type CrushNode struct {
	Selector Selector
}

type TestingNode struct {
	Children []Node
	CrushNode
	Weight float32
	Parent Node
	Failed bool
	ID     string
	Type   uint16
	Alg    string
}

func (n CrushNode) GetSelector() Selector {
	return n.Selector
}

func (n *CrushNode) SetSelector(Selector Selector) {
	n.Selector = Selector
}

func (n CrushNode) Select(input uint32, round uint32) Node {
	return n.GetSelector().Select(input, round)
}

func (n TestingNode) IsFailed() bool {
	return n.Failed
}

func (n TestingNode) IsLeaf() bool {
	return len(n.Children) == 0
}

func (n TestingNode) GetParent() Node {
	return n.Parent
}

func (n TestingNode) GetID() string {
	return n.ID
}

func (n TestingNode) GetWeight() float32 {
	return n.Weight
}

func (n TestingNode) GetType() uint16 {
	return n.Type
}

func (n TestingNode) GetChildren() []Node {
	return n.Children
}
