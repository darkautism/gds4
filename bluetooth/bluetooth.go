package bluetooth

import (
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"golang.org/x/sys/unix"
)

type BT struct {
	ctl_fd  int
	data_fd int
}

func (bt *BT) Read(p []byte) (n int, err error) {
	return unix.Read(bt.data_fd, p)
}

func (bt *BT) Write(p []byte) (n int, err error) {
	return unix.Write(bt.ctl_fd, p)
}

func (bt *BT) Close() error {
	err1 := unix.Close(bt.ctl_fd)
	err2 := unix.Close(bt.data_fd)
	if err1 != nil {
		return err1
	}
	return err2
}

func newL2Conn(addr [6]uint8, channel int) (int, error) {
	fd, err := unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_SEQPACKET, unix.BTPROTO_L2CAP)
	if err != nil {
		return -1, err
	}
	if err := unix.Connect(fd, &unix.SockaddrL2{
		PSM:  uint16(channel),
		Addr: addr,
	}); err != nil {
		return -1, err
	}
	return fd, nil
}

func NewBT(addrStr string) (io.ReadWriteCloser, error) {
	decoded, err := hex.DecodeString(strings.ReplaceAll(addrStr, ":", ""))
	if err != nil {
		return nil, fmt.Errorf("Input bluetooth address %s is invalid.", addrStr)
	}
	var addr [6]uint8
	var ret BT
	copy(addr[:], decoded)
	ret.ctl_fd, err = newL2Conn(addr, 0x11)
	if err != nil {
		return nil, err
	}
	ret.data_fd, err = newL2Conn(addr, 0x13)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
