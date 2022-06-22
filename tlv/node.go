package tlv

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type Node struct {
	//data []byte

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

	node = &Node{}

	offset := 0

	// Type
	start := offset
	d := data[offset]
	if (d & 0x1F) == 0x1F {
		for i := 0; i < len(data); i++ {
			d := data[offset]
			offset++
			if (d & 0x80) != 0x80 {
				break
			}
		}
	}
	node.Type = data[start:offset]

	// Length
	start = offset
	if data[offset]&0x80 != 0 {
		lengthLen := int(data[offset] & 0x7F)
		offset++
		for i := 0; i < lengthLen; i++ {
			node.Length = node.Length*256 + uint(data[offset])
			offset++
		}
	} else {
		node.Length = uint(data[offset])
		offset++
	}

	// Value
	start = offset
	if uint(offset)+node.Length > uint(len(data)) {
		return nil, errors.New("node parse overflow")
	}
	node.Value = data[start : offset+int(node.Length)]

	node.DataLen = uint(offset) + node.Length

	return node, nil
}

func (n Node) String() string {
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

func (n Node) IsConstructed() bool {
	/* If it's a constructed tag containing more tags
	 * check https://en.wikipedia.org/wiki/X.690#Encoding */
	return (n.Type[0] & (0x01 << 5)) != 0
}

func (n1 Node) Equal(n2 Node) bool {
	if !bytes.Equal(n1.Type, n2.Type) || n1.Length != n2.Length || !bytes.Equal(n1.Value, n2.Value) {
		return false
	}

	if len(n1.Nodes) != len(n2.Nodes) {
		return false
	}

	for i := 0; i < len(n1.Nodes); i++ {
		if !n1.Nodes[i].Equal(n2.Nodes[i]) {
			return false
		}
	}

	return true
}

func (ns Nodes) String() string {
	var sb strings.Builder

	for _, n := range ns {
		s := fmt.Sprintf("%s ", &n)
		sb.Write([]byte(s))
	}

	return sb.String()
}

func (ns1 Nodes) Equal(ns2 Nodes) bool {
	if len(ns1) != len(ns2) {
		return false
	}

	for i := 0; i < len(ns1); i++ {
		if !ns1[i].Equal(ns2[i]) {
			return false
		}
	}

	return true
}

func GetNodesWithType(nodes []Node, t []byte) []Node {
	ns := make([]Node, 0)

	for _, n := range nodes {
		if bytes.Equal(t, n.Type) {
			ns = append(ns, n)
		}
	}

	return ns
}
