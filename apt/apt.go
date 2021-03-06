package apt

import (
	"math"
	"math/rand"
	"reflect"
	"strconv"

	"github.com/LucasK1/gameswithgo/noise"
)

type Node interface {
	Eval(x, y float32) float32
	String() string
	SetParent(parent Node)
	GetParent() Node
	SetChildren([]Node)
	GetChildren() []Node
	AddRandom(node Node)
	AddLeaf(leaf Node) bool
	NodeCount() int
}

type BaseNode struct {
	Parent   Node
	Children []Node
}

func CopyTree(node Node, parent Node) Node {
	copy := reflect.New(reflect.ValueOf(node).Elem().Type()).Interface().(Node)

	copy.SetParent(parent)
	copyChildren := make([]Node, len(node.GetChildren()))
	copy.SetChildren(copyChildren)
	for i := range copyChildren {
		copyChildren[i] = CopyTree(node.GetChildren()[i], copy)
	}
	return copy
}

func ReplaceNode(old Node, new Node) {
	oldParent := old.GetParent()
	if oldParent != nil {
		for i, child := range oldParent.GetChildren() {
			if child == old {
				oldParent.GetChildren()[i] = new
			}
		}
	}
	new.SetParent(oldParent)
}

func GetNthNode(node Node, n, count int) (Node, int) {
	if n == count {
		return node, count
	}
	var result Node
	for _, child := range node.GetChildren() {
		count++
		result, count = GetNthNode(child, n, count)
		if result != nil {
			return result, count
		}
	}
	return nil, count
}

func Mutate(node Node) Node {
	r := rand.Intn(23)
	var mutatedNode Node

	if r <= 19 {
		mutatedNode = GetRandomNode()
	} else {
		mutatedNode = GetRandomLeaf()
	}

	// Fix up parent's child pointer to point to the new node
	if node.GetParent() != nil {
		for i, parentChild := range node.GetParent().GetChildren() {
			if parentChild == node {
				node.GetParent().GetChildren()[i] = mutatedNode
			}
		}
	}

	// Add children from the old node to the new mutated node
	for i, child := range node.GetChildren() {
		if i >= len(mutatedNode.GetChildren()) {
			break
		}
		mutatedNode.GetChildren()[i] = child
		child.SetParent(mutatedNode)
	}

	// Any nil children get filled in with a random leaf
	for i, child := range mutatedNode.GetChildren() {
		if child == nil {
			leaf := GetRandomLeaf()
			leaf.SetParent(mutatedNode)
			mutatedNode.GetChildren()[i] = leaf
		}
	}
	mutatedNode.SetParent(node.GetParent())
	return mutatedNode
}

func (node *BaseNode) Eval(x, y float32) float32 {
	panic("Tried to call eval on basenode")
}

func (node *BaseNode) String() string {
	panic("Tried to call String on basenode")
}

func (node *BaseNode) SetParent(parent Node) {
	node.Parent = parent
}

func (node *BaseNode) GetParent() Node {
	return node.Parent
}

func (node *BaseNode) SetChildren(children []Node) {
	node.Children = children
}

func (node *BaseNode) GetChildren() []Node {
	return node.Children
}

func (node *BaseNode) AddRandom(nodeToAdd Node) {
	addIndex := rand.Intn(len(node.Children))
	if node.Children[addIndex] == nil {
		nodeToAdd.SetParent(node)
		node.Children[addIndex] = nodeToAdd
	} else {
		node.Children[addIndex].AddRandom(nodeToAdd)
	}
}

func (node *BaseNode) AddLeaf(leaf Node) bool {
	for i, child := range node.Children {
		if child == nil {
			leaf.SetParent(node)
			node.Children[i] = leaf
			return true
		} else if node.Children[i].AddLeaf(leaf) {
			return true
		}
	}
	return false
}

func (node *BaseNode) NodeCount() int {
	count := 1
	for _, child := range node.Children {
		count += child.NodeCount()
	}
	return count
}

type OpSin struct {
	BaseNode
}

func NewOpSin() *OpSin {
	return &OpSin{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpSin) Eval(x, y float32) float32 {
	return float32(math.Sin(float64(op.Children[0].Eval(x, y))))
}

func (op *OpSin) String() string {
	return "( Sin " + op.Children[0].String() + " )"
}

type OpCos struct {
	BaseNode
}

func NewOpCos() *OpCos {
	return &OpCos{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpCos) Eval(x, y float32) float32 {
	return float32(math.Cos(float64(op.Children[0].Eval(x, y))))
}

func (op *OpCos) String() string {
	return "( Cos " + op.Children[0].String() + " )"
}

type OpTan struct {
	BaseNode
}

func NewOpTan() *OpTan {
	return &OpTan{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpTan) Eval(x, y float32) float32 {
	return float32(math.Tan(float64(op.Children[0].Eval(x, y))))
}

type OpAtan struct {
	BaseNode
}

func NewOpAtan() *OpAtan {
	return &OpAtan{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpAtan) Eval(x, y float32) float32 {
	return float32(math.Atan(float64(op.Children[0].Eval(x, y))))
}

func (op *OpAtan) String() string {
	return "( Atan " + op.Children[0].String() + " )"
}

type OpLog2 struct {
	BaseNode
}

func NewOpLog2() *OpLog2 {
	return &OpLog2{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpLog2) Eval(x, y float32) float32 {
	return float32(math.Log2(float64(op.Children[0].Eval(x, y))))
}

func (op *OpLog2) String() string {
	return "( Log2 " + op.Children[0].String() + " )"
}

type OpAbs struct {
	BaseNode
}

func NewOpAbs() *OpAbs {
	return &OpAbs{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpAbs) Eval(x, y float32) float32 {
	return float32(math.Abs(float64(op.Children[0].Eval(x, y))))
}

func (op *OpAbs) String() string {
	return "( Abs " + op.Children[0].String() + " )"
}

type OpCeil struct {
	BaseNode
}

func NewOpCeil() *OpCeil {
	return &OpCeil{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpCeil) Eval(x, y float32) float32 {
	return float32(math.Ceil(float64(op.Children[0].Eval(x, y))))
}

func (op *OpCeil) String() string {
	return "( Ceil " + op.Children[0].String() + " )"
}

type OpFloor struct {
	BaseNode
}

func NewOpFloor() *OpFloor {
	return &OpFloor{BaseNode{nil, make([]Node, 1)}}
}

func (op *OpFloor) Eval(x, y float32) float32 {
	return float32(math.Floor(float64(op.Children[0].Eval(x, y))))
}

func (op *OpFloor) String() string {
	return "( Floor " + op.Children[0].String() + " )"
}

type OpAtan2 struct {
	BaseNode
}

func NewOpAtan2() *OpAtan2 {
	return &OpAtan2{BaseNode{nil, make([]Node, 2)}}
}

func (op *OpAtan2) Eval(x, y float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func (op *OpAtan2) String() string {
	return "( Atan2 " + op.Children[0].String() + " " + op.Children[1].String() + " )"
}

type OpNoise struct {
	BaseNode
}

func NewOpNoise() *OpNoise {
	return &OpNoise{BaseNode{nil, make([]Node, 2)}}
}

func (op *OpNoise) Eval(x, y float32) float32 {
	return 80*noise.Snoise2(op.Children[0].Eval(x, y), op.Children[1].Eval(x, y)) - 2
}

func (op *OpNoise) String() string {
	return "( SimplexNoise " + op.Children[0].String() + " " + op.Children[1].String() + " )"
}

type OpTurbulence struct {
	BaseNode
}

func NewOpTurbulence() *OpTurbulence {
	return &OpTurbulence{BaseNode{nil, make([]Node, 3)}}
}

func (op *OpTurbulence) Eval(x, y float32) float32 {
	return 2*6.96*noise.Turbulence(op.Children[0].Eval(x, y), op.Children[2].Eval(x, y), 5*op.Children[1].Eval(x, y), 0.5, 2, 3) - 1
}

func (op *OpTurbulence) String() string {
	return "( Turbulence " + op.Children[0].String() + " " + op.Children[2].String() + " " + op.Children[1].String() + " )"
}

type OpFBM struct {
	BaseNode
}

func NewOpFBM() *OpFBM {
	return &OpFBM{BaseNode{nil, make([]Node, 3)}}
}

func (op *OpFBM) Eval(x, y float32) float32 {
	return 2*3.627*noise.Fbm2(op.Children[0].Eval(x, y), op.Children[2].Eval(x, y), 5*op.Children[1].Eval(x, y), 0.5, 2, 3) + 0.492 - 1
}

func (op *OpFBM) String() string {
	return "( FBM " + op.Children[0].String() + " " + op.Children[2].String() + " " + op.Children[1].String() + " )"
}

type OpLerp struct {
	BaseNode
}

func NewOpLerp() *OpLerp {
	return &OpLerp{BaseNode{nil, make([]Node, 3)}}
}

func (op *OpLerp) Eval(x, y float32) float32 {
	return op.Children[0].Eval(x, y) + op.Children[2].Eval(x, y)*(op.Children[1].Eval(x, y)-op.Children[0].Eval(x, y))
}

func (op *OpLerp) String() string {
	return "( Lerp " + op.Children[0].String() + " " + op.Children[2].String() + " " + op.Children[1].String() + " )"
}

type OpPlus struct {
	BaseNode
}

func NewOpPlus() *OpPlus {
	return &OpPlus{BaseNode{nil, make([]Node, 2)}}
}

func (op *OpPlus) Eval(x, y float32) float32 {
	return op.Children[0].Eval(x, y) + op.Children[1].Eval(x, y)
}

func (op *OpPlus) String() string {
	return "( + " + op.Children[0].String() + " " + op.Children[1].String() + " )"
}

type OpMinus struct {
	BaseNode
}

func NewOpMinus() *OpMinus {
	return &OpMinus{BaseNode{nil, make([]Node, 2)}}
}

func (op *OpMinus) Eval(x, y float32) float32 {
	return op.Children[0].Eval(x, y) - op.Children[1].Eval(x, y)
}

func (op *OpMinus) String() string {
	return "( - " + op.Children[0].String() + " " + op.Children[1].String() + " )"
}

type OpMult struct {
	BaseNode
}

func NewOpMult() *OpMult {
	return &OpMult{BaseNode{nil, make([]Node, 2)}}
}

func (op *OpMult) Eval(x, y float32) float32 {
	return op.Children[0].Eval(x, y) * op.Children[1].Eval(x, y)
}

func (op *OpMult) String() string {
	return "( * " + op.Children[0].String() + " " + op.Children[1].String() + " )"
}

type OpDiv struct {
	BaseNode
}

func NewOpDiv() *OpDiv {
	return &OpDiv{BaseNode{nil, make([]Node, 2)}}
}

func (op *OpDiv) Eval(x, y float32) float32 {
	return op.Children[0].Eval(x, y) / op.Children[1].Eval(x, y)
}

func (op *OpDiv) String() string {
	return "( / " + op.Children[0].String() + " " + op.Children[1].String() + " )"
}

type OpX struct {
	BaseNode
}

func NewOpX() *OpX {
	return &OpX{BaseNode{nil, make([]Node, 0)}}
}

func (op *OpX) Eval(x, y float32) float32 {
	return x
}

func (op *OpX) String() string {
	return "X"
}

type OpY struct {
	BaseNode
}

func NewOpY() *OpY {
	return &OpY{BaseNode{nil, make([]Node, 0)}}
}

func (op *OpY) Eval(x, y float32) float32 {
	return y
}

func (op *OpY) String() string {
	return "Y"
}

type OpConstant struct {
	BaseNode
	value float32
}

func NewOpConstant() *OpConstant {
	return &OpConstant{BaseNode{nil, make([]Node, 0)}, rand.Float32()*2 - 1}
}

func (op *OpConstant) Eval(x, y float32) float32 {
	return op.value
}

func (op *OpConstant) String() string {
	return strconv.FormatFloat(float64(op.value), 'f', 9, 32)
}

func GetRandomNode() Node {
	r := rand.Intn(16)

	switch r {
	case 0:
		return NewOpPlus()
	case 1:
		return NewOpMinus()
	case 2:
		return NewOpMult()
	case 3:
		return NewOpDiv()
	case 4:
		return NewOpAtan2()
	case 5:
		return NewOpAtan()
	case 6:
		return NewOpCos()
	case 7:
		return NewOpSin()
	case 8:
		return NewOpNoise()
	case 9:
		return NewOpLerp()
	case 10:
		return NewOpLog2()
	case 11:
		return NewOpAbs()
	case 12:
		return NewOpCeil()
	case 13:
		return NewOpFloor()
	case 14:
		return NewOpTurbulence()
	case 15:
		return NewOpFBM()
	}
	panic("Get Random Node failed")
}

func GetRandomLeaf() Node {
	r := rand.Intn(3)

	switch r {
	case 0:
		return NewOpX()
	case 1:
		return NewOpY()
	case 2:
		return NewOpConstant()
	}
	panic("Get Random Leaf Node failed")
}
