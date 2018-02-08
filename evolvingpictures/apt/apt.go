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

type TripleNode struct {
	LeftChild   Node
	MiddleChild Node
	RightChild  Node
}

func (triple *TripleNode) NodeCounts() (nodeCount, nilCount int) {
	var leftCount, leftNilCount, middleCount, middleNilCount, rightCount, rightNilCount int
	if triple.LeftChild == nil {
		leftNilCount = 1
		leftCount = 0
	} else {
		leftCount, leftNilCount = triple.LeftChild.NodeCounts()
	}

	if triple.MiddleChild == nil {
		middleNilCount = 1
		middleCount = 0
	} else {
		middleCount, middleNilCount = triple.MiddleChild.NodeCounts()
	}

	if triple.RightChild == nil {
		rightNilCount = 1
		rightCount = 0
	} else {
		rightCount, rightNilCount = triple.RightChild.NodeCounts()
	}

	return 1 + leftCount + middleCount + rightCount, leftNilCount + middleNilCount + rightNilCount
}

func (triple *TripleNode) AddRandom(node Node) {
	r := rand.Intn(3)
	if r == 0 {
		if triple.LeftChild == nil {
			triple.LeftChild = node
		} else {
			triple.LeftChild.AddRandom(node)
		}
	} else if r == 1 {
		if triple.MiddleChild == nil {
			triple.MiddleChild = node
		} else {
			triple.MiddleChild.AddRandom(node)
		}
	} else {
		if triple.RightChild == nil {
			triple.RightChild = node
		} else {
			triple.RightChild.AddRandom(node)
		}
	}
}

type OpLerp struct {
	TripleNode
}

func (op *OpLerp) Eval(x, y float32) float32 {

	a := op.LeftChild.Eval(x, y)
	b := op.MiddleChild.Eval(x, y)
	pct := op.RightChild.Eval(x, y)
	return a + pct*(b-a)
}

func (op *OpLerp) String() string {
	return "( Lerp " + op.LeftChild.String() + " " + op.MiddleChild.String() + " " + op.RightChild.String() + " )"
}

type OpClip struct {
	DoubleNode
}

func (op *OpClip) Eval(x, y float32) float32 {
	value := op.LeftChild.Eval(x, y)
	max := float32(math.Abs(float64(op.RightChild.Eval(x, y))))
	if value > max {
		return max
	} else if value < -max {
		return -max
	}
	return value
}

func (op *OpClip) String() string {
	return "( Clip " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpWrap struct {
	SingleNode
}

func (op *OpWrap) Eval(x, y float32) float32 {
	f := op.Child.Eval(x, y)
	temp := (f - -1.0) / (2.0)
	return -1.0 + 2.0*(temp-float32(math.Floor(float64(temp))))
}
func (op *OpWrap) String() string {
	return "(Wrap " + op.Child.String() + ")"
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

type OpLog2 struct {
	SingleNode
}

func (op *OpLog2) Eval(x, y float32) float32 {
	return float32(math.Log2(float64(op.Child.Eval(x, y))))
}

func (op *OpLog2) String() string {
	return "( Log2 " + op.Child.String() + " )"
}

type OpSquare struct {
	SingleNode
}

func (op *OpSquare) Eval(x, y float32) float32 {
	value := op.Child.Eval(x, y)
	return value * value
}

func (op *OpSquare) String() string {
	return "( Square " + op.Child.String() + " )"
}

type OpNegate struct {
	SingleNode
}

func (op *OpNegate) Eval(x, y float32) float32 {
	return -op.Child.Eval(x, y)
}

func (op *OpNegate) String() string {
	return "( Negate " + op.Child.String() + " )"
}

type OpCeil struct {
	SingleNode
}

func (op *OpCeil) Eval(x, y float32) float32 {
	return float32(math.Ceil(float64(op.Child.Eval(x, y))))
}

func (op *OpCeil) String() string {
	return "( Ceil " + op.Child.String() + " )"
}

type OpFloor struct {
	SingleNode
}

func (op *OpFloor) Eval(x, y float32) float32 {
	return float32(math.Floor(float64(op.Child.Eval(x, y))))
}

func (op *OpFloor) String() string {
	return "( Floor " + op.Child.String() + " )"
}

type OpAbs struct {
	SingleNode
}

func (op *OpAbs) Eval(x, y float32) float32 {
	return float32(math.Abs(float64(op.Child.Eval(x, y))))
}

func (op *OpAbs) String() string {
	return "( Abs " + op.Child.String() + " )"
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

type OpFBM struct {
	TripleNode
}

func (op *OpFBM) Eval(x, y float32) float32 {
	return 2*3.627*noise.Fbm2(op.LeftChild.Eval(x, y), op.RightChild.Eval(x, y), 5*op.MiddleChild.Eval(x, y), 0.5, 2, 3) + .492 - 1
}

func (op *OpFBM) String() string {
	return "( FBM " + op.LeftChild.String() + " " + op.RightChild.String() + " " + op.MiddleChild.String() + " )"
}

type OpTurbulence struct {
	TripleNode
}

func (op *OpTurbulence) Eval(x, y float32) float32 {
	return 2*6.96*noise.Turbulence(op.LeftChild.Eval(x, y), op.RightChild.Eval(x, y), 5*op.MiddleChild.Eval(x, y), 0.5, 2, 3) - 1
}

func (op *OpTurbulence) String() string {
	return "( Turbulence " + op.LeftChild.String() + " " + op.RightChild.String() + " " + op.MiddleChild.String() + " )"
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
	r := rand.Intn(20)
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
	case 9:
		return &OpSquare{}
	case 10:
		return &OpLog2{}
	case 11:
		return &OpNegate{}
	case 12:
		return &OpCeil{}
	case 13:
		return &OpFloor{}
	case 14:
		return &OpLerp{}
	case 15:
		return &OpAbs{}
	case 16:
		return &OpClip{}
	case 17:
		return &OpWrap{}
	case 18:
		return &OpFBM{}
	case 19:
		return &OpTurbulence{}
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
