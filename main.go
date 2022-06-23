package main

import (
	"github.com/waynewu411/gopkg/logger"
	"github.com/waynewu411/gopkg/tlv"
)

func main() {
	logger.InitLogger()
	defer logger.Sync()

	logger.Infof("gopkg")

	data := []byte("\xFF\x03\x25\xDF\x02\x22\x39\x39\x36\x30\x30\x30\x30\x30\x30\x30\x30\x31\x30\x32\x30\x30\x5F\x54\x5A\x4D\x4F\x43\x4B\x55\x50\x5F\x46\x49\x4C\x45\x2E\x50\x33\x54")
	nodes, err := tlv.DecodeWithData(data)
	if nil != err {
		logger.Errorf(err, "tlv decode failed")
		return
	}
	logger.Debugf("tlv parsed, nodes = %v", nodes)
}
