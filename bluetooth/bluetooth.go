package bluetooth

import (
	"encoding/hex"
	"fmt"
	"strings"
)

struct BT {
	int ctl_fd;
	int data_fd;
}

func (bt *BT) Read(p []byte) (n int, err error) {
	return unix.Read(bt.data_fd, status_b)
}

func (bt *BT) Write(p []byte) (n int, err error) {
	return unix.Write(bt.ctl_fd, p);
}

func (bt *BT) Close() error {
	unix.Close(bt.Ctrl)
	unix.Close(bt.Data)
}

func NewBT(addrStr string) (io.ReadWriteCloser) {
	decoded, err := hex.DecodeString(strings.ReplaceAll(addrStr, ":", ""))
	if err != nil {
		return nil, fmt.Errorf("Input bluetooth address %s is invalid.", addrStr)
	}
	var addr [6]uint8
	var ret BT;
	copy(addr[:], decoded)
	ret.ctl_fd, err := newL2Conn(addr, 0x11)
	if err != nil {
		return nil, err
	}
	ret.data_fd, err := newL2Conn(addr, 0x13)
	if err != nil {
		return nil, err
	}
	return &ret;
}