package tlv

import "errors"

func DecodeWithData(data []byte) (nodes Nodes, err error) {
	offset := 0
	for {
		if offset == len(data) {
			break
		}

		node, err := DecodeNodeWithData(data[offset:])
		if nil != err {
			return nil, err
		}

		// Check container tag
		if node.IsContainer() {
			ns, err := DecodeWithData(node.Value)
			if nil != err {
				return nil, err
			}
			node.Nodes = ns
		}

		offset += int(node.DataLen)

		if offset > len(data) {
			return nil, errors.New("decoder overflow")
		}

		nodes = append(nodes, *node)
	}

	return nodes, nil
}
