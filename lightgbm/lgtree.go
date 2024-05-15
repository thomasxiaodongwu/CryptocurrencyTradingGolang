package lightgbm

import (
	"math"
)

const (
	categorical = 1 << 0
	defaultLeft = 1 << 1
	leftLeaf    = 1 << 2
	rightLeaf   = 1 << 3
	missingZero = 1 << 4
	missingNan  = 1 << 5
	catOneHot   = 1 << 6
	catSmall    = 1 << 7
)

const zeroThreshold = 1e-35

type lgNode struct {
	Threshold float64
	Left      uint32
	Right     uint32
	Feature   uint32
	Flags     uint8
}

type lgTree struct {
	nodes         []lgNode
	leafValues    []float64
	catBoundaries []uint32
	catThresholds []uint32
	nCategorical  uint32

	// xwu
	Pred  float64
	Fvals []float64
}

func (t *lgTree) numericalDecision(node *lgNode, fval float64) bool {
	isnan := math.IsNaN(fval)
	fmn := node.Flags & missingNan
	if isnan && (fmn == 0) {
		fval = 0.0
	}
	if ((node.Flags&missingZero > 0) && isZero(fval)) || ((fmn > 0) && isnan) {
		return node.Flags&defaultLeft > 0
	}
	// Note: LightGBM uses `<=`, but XGBoost uses `<`
	return fval <= node.Threshold
}

func (t *lgTree) categoricalDecision(node *lgNode, fval float64) bool {
	ifval := int32(fval)
	if ifval < 0 {
		return false
	} else if math.IsNaN(fval) {
		if node.Flags&missingNan > 0 {
			return false
		}
		ifval = 0
	}
	if node.Flags&catOneHot > 0 {
		return int32(node.Threshold) == ifval
	} else if node.Flags&catSmall > 0 {
		return FindInBitsetUint32(uint32(node.Threshold), uint32(ifval))
	}
	return t.findInBitset(uint32(node.Threshold), uint32(ifval))
}

func (t *lgTree) decision(node *lgNode, fval float64) bool {
	if node.Flags&categorical > 0 { // 如果是类别的特征
		ifval := int32(fval)
		if ifval < 0 {
			return false
		} else if math.IsNaN(fval) {
			if node.Flags&missingNan > 0 {
				return false
			}
			ifval = 0
		}
		if node.Flags&catOneHot > 0 {
			return int32(node.Threshold) == ifval
		} else if node.Flags&catSmall > 0 {
			return FindInBitsetUint32(uint32(node.Threshold), uint32(ifval))
		}
		return t.findInBitset(uint32(node.Threshold), uint32(ifval))
	}

	isnan := math.IsNaN(fval)
	fmn := node.Flags & missingNan
	if isnan && (fmn == 0) {
		fval = 0.0
	}
	if ((node.Flags&missingZero > 0) && isZero(fval)) || ((fmn > 0) && isnan) {
		return node.Flags&defaultLeft > 0
	}
	// Note: LightGBM uses `<=`, but XGBoost uses `<`
	return fval <= node.Threshold
}

// func (t *lgTree) decision(node *lgNode, fval float64) bool {
// 	if node.Flags&categorical > 0 { // 如果是类别的特征
// 		return t.categoricalDecision(node, fval)
// 	}
// 	return t.numericalDecision(node, fval)
// }

func (t *lgTree) predict(fvals []float64) (float64, uint32) {
	// TODO:输入和输出没有搞懂暂时
	if len(t.nodes) == 0 { // 这里应该如果是一个空的树，就返回左子节点的值
		return t.leafValues[0], 0
	}
	idx := uint32(0) // 从根节点开始搜索
	for {
		node := &t.nodes[idx]
		left := t.decision(node, fvals[node.Feature]) // 决定是往左边走还是右边走
		if left {
			if node.Flags&leftLeaf > 0 { // 如果是叶子节点，那么就返回结果，这里应该是位运算
				t.Pred = t.leafValues[node.Left]
				return t.leafValues[node.Left], node.Left
			}
			idx = node.Left // 这里左边和右边不同方式的更新应该是树的存储方式导致的
		} else {
			if node.Flags&rightLeaf > 0 {
				t.Pred = t.leafValues[node.Right]
				return t.leafValues[node.Right], node.Right
			}
			idx++
		}
	}
}

func (t *lgTree) findInBitset(idx uint32, pos uint32) bool {
	i1 := pos / 32
	idxS := t.catBoundaries[idx]
	idxE := t.catBoundaries[idx+1]
	if i1 >= (idxE - idxS) {
		return false
	}
	i2 := pos % 32
	return (t.catThresholds[idxS+i1]>>i2)&1 > 0
}

func (t *lgTree) nLeaves() int {
	return len(t.nodes) + 1
}

func (t *lgTree) nNodes() int {
	return len(t.nodes)
}

func isZero(fval float64) bool {
	return (fval > -zeroThreshold && fval <= zeroThreshold)
}

func categoricalNode(feature uint32, missingType uint8, threshold uint32, catType uint8) lgNode {
	node := lgNode{}
	node.Feature = feature
	node.Flags = categorical | missingType | catType
	node.Threshold = float64(threshold)
	return node
}

func numericalNode(feature uint32, missingType uint8, threshold float64, defaultType uint8) lgNode {
	node := lgNode{}
	node.Feature = feature
	node.Flags = missingType | defaultType
	node.Threshold = threshold
	return node
}
