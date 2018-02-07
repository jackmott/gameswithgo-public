//TODO backtracking on addrandom

package apt

import (
	"fmt"
	"github.com/jackmott/noise"
	"math"
	"math/rand"
	"strconv"
)

// + / * - Sin Cos Atan SimplexNoise X Y Constants...
// Leaf Node (0 children)
// Single Node (sin/cos)
// DoubleNode (+, -)

type Node interface {
	Eval(x, y float32) float32
	String() string
	AddRandom(node Node)
	NodeCounts() (nodeCount, nilCount int)
}

type LeafNode struct{}

func (leaf *LeafNode) AddRandom(node Node) {
	//panic("ERROR: You tried to add a node to a leaf node")
	fmt.Println("tried to add a leaf node")
}

func (leaf *LeafNode) NodeCounts() (nodeCount, nilCount int) {
	return 1, 0
}

type SingleNode struct {
	Child Node
}

func (single *SingleNode) AddRandom(node Node) {
	if single.Child == nil {
		single.Child = node
	} else {
		single.Child.AddRandom(node)
	}
}

func (single *SingleNode) NodeCounts() (nodeCount, nilCount int) {
	if single.Child == nil {
		return 1, 1
	} else {
		childNodeCount, childNilCount := single.Child.NodeCounts()
		return 1 + childNodeCount, childNilCount
	}
}

type DoubleNode struct {
	LeftChild  Node
	RightChild Node
}

func (double *DoubleNode) NodeCounts() (nodeCount, nilCount int) {
	var leftCount, leftNilCount, rightCount, rightNilCount int
	if double.LeftChild == nil {
		leftNilCount = 1
		leftCount = 0
	} else {
		leftCount, leftNilCount = double.LeftChild.NodeCounts()
	}

	if double.RightChild == nil {
		rightNilCount = 1
		rightCount = 0
	} else {
		rightCount, rightNilCount = double.RightChild.NodeCounts()
	}

	return 1 + leftCount + rightCount, leftNilCount + rightNilCount
}

func (double *DoubleNode) AddRandom(node Node) {
	r := rand.Intn(2)
	if r == 0 {
		if double.LeftChild == nil {
			double.LeftChild = node
		} else {
			double.LeftChild.AddRandom(node)
		}
	} else {
		if double.RightChild == nil {
			double.RightChild = node
		} else {
			double.RightChild.AddRandom(node)
		}
	}
}

type OpSin struct {
	SingleNode
}

func (op *OpSin) Eval(x, y float32) float32 {
	return float32(math.Sin(float64(op.Child.Eval(x, y))))
}

func (op *OpSin) String() string {
	return "( Sin " + op.Child.String() + " )"
}

type OpCos struct {
	SingleNode
}

func (op *OpCos) Eval(x, y float32) float32 {
	return float32(math.Cos(float64(op.Child.Eval(x, y))))
}

func (op *OpCos) String() string {
	return "( Cos " + op.Child.String() + " )"
}

type OpAtan struct {
	SingleNode
}

func (op *OpAtan) Eval(x, y float32) float32 {
	return float32(math.Atan(float64(op.Child.Eval(x, y))))
}

func (op *OpAtan) String() string {
	return "( Atan " + op.Child.String() + " )"
}

type OpNoise struct {
	DoubleNode
}

func (op *OpNoise) Eval(x, y float32) float32 {
	return 80*noise.Snoise2(op.LeftChild.Eval(x, y), op.RightChild.Eval(x, y)) - 2.0
}

func (op *OpNoise) String() string {
	return "( SimplexNoise " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpPlus struct {
	DoubleNode
}

func (op *OpPlus) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) + op.RightChild.Eval(x, y)
}

func (op *OpPlus) String() string {
	return "( + " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpMinus struct {
	DoubleNode
}

func (op *OpMinus) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) - op.RightChild.Eval(x, y)
}

func (op *OpMinus) String() string {
	return "( - " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpMult struct {
	DoubleNode
}

func (op *OpMult) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) * op.RightChild.Eval(x, y)
}

func (op *OpMult) String() string {
	return "( * " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpDiv struct {
	DoubleNode
}

func (op *OpDiv) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) / op.RightChild.Eval(x, y)
}

func (op *OpDiv) String() string {
	return "( / " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpAtan2 struct {
	DoubleNode
}

func (op *OpAtan2) Eval(x, y float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func (op *OpAtan2) String() string {
	return "( Atan2 " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpX struct {
	LeafNode
}

func (op *OpX) Eval(x, y float32) float32 {
	return x
}

func (op *OpX) String() string {
	return "X"
}

type OpY struct {
	LeafNode
}

func (op *OpY) Eval(x, y float32) float32 {
	return y
}

func (op *OpY) String() string {
	return "Y"
}

type OpConstant struct {
	LeafNode
	value float32
}

func (op *OpConstant) Eval(x, y float32) float32 {
	return op.value
}

func (op *OpConstant) String() string {
	return strconv.FormatFloat(float64(op.value), 'f', 9, 32)
}

func GetRandomNode() Node {
	r := rand.Intn(9)
	switch r {
	case 0:
		return &OpPlus{}
	case 1:
		return &OpMinus{}
	case 2:
		return &OpMult{}
	case 3:
		return &OpDiv{}
	case 4:
		return &OpAtan2{}
	case 5:
		return &OpAtan{}
	case 6:
		return &OpCos{}
	case 7:
		return &OpSin{}
	case 8:
		return &OpNoise{}
	}
	panic("Get Random Noise Failed")
}

func GetRandomLeaf() Node {
	r := rand.Intn(3)
	switch r {
	case 0:
		return &OpX{}
	case 1:
		return &OpY{}
	case 2:
		return &OpConstant{LeafNode{}, rand.Float32()*2 - 1}
	}
	panic("Error in get random Leaf")
}
