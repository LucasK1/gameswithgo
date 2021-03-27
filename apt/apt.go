package apt

import (
	"math"
	"strconv"

	"github.com/LucasK1/gameswithgo/noise"
)

type Node interface {
	Eval(x, y float32) float32
	String() string
}

type LeafNode struct{}

type SingleNode struct {
	Child Node
}

type DoubleNode struct {
	LeftChild  Node
	RightChild Node
}

type OpSin SingleNode

func (op *OpSin) Eval(x, y float32) float32 {
	return float32(math.Sin(float64(op.Child.Eval(x, y))))
}

func (op *OpSin) String() string {
	return "( Sin " + op.Child.String() + " )"
}

type OpCos SingleNode

func (op *OpCos) Eval(x, y float32) float32 {
	return float32(math.Cos(float64(op.Child.Eval(x, y))))
}

func (op *OpCos) String() string {
	return "( Cos " + op.Child.String() + " )"
}

type OpAtan SingleNode

func (op *OpAtan) Eval(x, y float32) float32 {
	return float32(math.Atan(float64(op.Child.Eval(x, y))))
}

func (op *OpAtan) String() string {
	return "( Atan " + op.Child.String() + " )"
}

type OpAtan2 DoubleNode

func (op *OpAtan2) Eval(x, y float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func (op *OpAtan2) String() string {
	return "( Atan2 " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpNoise DoubleNode

func (op *OpNoise) Eval(x, y float32) float32 {
	return 80*noise.Snoise2(op.LeftChild.Eval(x, y), op.RightChild.Eval(x, y)) - 2
}

func (op *OpNoise) String() string {
	return "( SimplexNoise " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpPlus DoubleNode

func (op *OpPlus) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) + op.RightChild.Eval(x, y)
}

func (op *OpPlus) String() string {
	return "( + " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpMinus DoubleNode

func (op *OpMinus) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) - op.RightChild.Eval(x, y)
}

func (op *OpMinus) String() string {
	return "( - " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpMult DoubleNode

func (op *OpMult) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) * op.RightChild.Eval(x, y)
}

func (op *OpMult) String() string {
	return "( * " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpDiv DoubleNode

func (op *OpDiv) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) / op.RightChild.Eval(x, y)
}

func (op *OpDiv) String() string {
	return "( / " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OpX LeafNode

func (op *OpX) Eval(x, y float32) float32 {
	return x
}

func (op *OpX) String() string {
	return "X"
}

type OpY LeafNode

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
