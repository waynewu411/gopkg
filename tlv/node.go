package tlv

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type Node struct {
	data []byte

	Type    []byte
	Length  uint
	Value   []byte
	DataLen uint // Total Length of the whole TLV

	Nodes Nodes
}

type Nodes []Node

func DecodeNodeWithData(data []byte) (node *Node, err error) {
	if len(data) == 0 {
		return nil, errors.New("data len zero")
	}

	node = &Node{data: data}

	offset := 0

	// Type
	start := offset
	d := node.data[offset]
	if 0x1F == (d & 0x1F) {
		for i := 0; i < len(node.data); i++ {
			d := node.data[offset]
			offset++
			if 0x80 != (d & 0x80) {
				break
			}
		}
	}
	node.Type = node.data[start:offset]

	// Length
	start = offset
	if node.data[offset]&0x80 != 0 {
		lengthLen := int(node.data[offset] & 0x7F)
		offset++
		for i := 0; i < lengthLen; i++ {
			node.Length = node.Length*256 + uint(node.data[offset])
			offset++
		}
	} else {
		node.Length = uint(node.data[offset])
		offset++
	}

	// Value
	start = offset
	if uint(offset)+node.Length > uint(len(node.data)) {
		return nil, errors.New("node parse overflow")
	}
	node.Value = node.data[start : offset+int(node.Length)]

	node.DataLen = uint(offset) + node.Length

	return node, nil
}

func (n *Node) String() string {
	var sb strings.Builder
	s := fmt.Sprintf("[%s/%d: %s", strings.ToUpper(hex.EncodeToString(n.Type)), n.Length, strings.ToUpper(hex.EncodeToString(n.Value)))
	sb.Write([]byte(s))
	for _, node := range n.Nodes {
		s = fmt.Sprintf("  %s", &node)
		sb.Write([]byte(s))
	}
	sb.Write([]byte("]"))

	return sb.String()
}

var containerTypes = [][]byte{
	[]byte("\xFF\x01"),
	[]byte("\xFF\x03")}

func (n *Node) IsContainer() bool {
	for _, Type := range containerTypes {
		if 0 == bytes.Compare(Type, n.Type) {
			return true
		}
	}

	return false
}

func (ns Nodes) String() string {
	var sb strings.Builder

	for _, n := range ns {
		s := fmt.Sprintf("%s ", &n)
		sb.Write([]byte(s))
	}

	return sb.String()
}

func GetNodesWithType(nodes []Node, t []byte) []Node {
	ns := make([]Node, 0)

	for _, n := range nodes {
		if 0 == bytes.Compare(t, n.Type) {
			ns = append(ns, n)
		}
	}

	return ns
}
