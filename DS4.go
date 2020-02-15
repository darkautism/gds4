package gds4

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"image/color"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func BTAddrString2Addr(addr string) (*[6]uint8, error) {
	decoded, err := hex.DecodeString(strings.ReplaceAll(addr, ":", ""))
	if err != nil {
		return nil, err
	}
	var ret [6]uint8
	copy(ret[:], decoded)
	return &ret, nil
}

func (dp DS4_Packet) Up() bool {
	return dp.PAD|0xF == 0 || dp.PAD|0xF == 1 || dp.PAD|0xF == 7
}
func (dp DS4_Packet) Right() bool {
	return dp.PAD|0xF == 1 || dp.PAD|0xF == 2 || dp.PAD|0xF == 3
}
func (dp DS4_Packet) Down() bool {
	return dp.PAD|0xF == 3 || dp.PAD|0xF == 4 || dp.PAD|0xF == 5
}
func (dp DS4_Packet) Left() bool {
	return dp.PAD|0xF == 5 || dp.PAD|0xF == 6 || dp.PAD|0xF == 7
}
func (dp DS4_Packet) Square() bool {
	return dp.PAD&0b10000 == 1
}
func (dp DS4_Packet) X() bool {
	return dp.PAD&0b100000 == 1
}
func (dp DS4_Packet) O() bool {
	return dp.PAD&0b1000000 == 1
}
func (dp DS4_Packet) Triangle() bool {
	return dp.PAD&0b10000000 == 1
}
func (dp DS4_Packet) L1() bool {
	return dp.BTN&0b0001 == 1
}
func (dp DS4_Packet) R1() bool {
	return dp.BTN&0b0010 == 1
}
func (dp DS4_Packet) L2() bool {
	return dp.BTN&0b0100 == 1
}
func (dp DS4_Packet) R2() bool {
	return dp.BTN&0b1000 == 1
}
func (dp DS4_Packet) Share() bool {
	return dp.BTN&0b10000 == 1
}
func (dp DS4_Packet) Option() bool {
	return dp.BTN&0b100000 == 1
}
func (dp DS4_Packet) L3() bool {
	return dp.BTN&0b1000000 == 1
}
func (dp DS4_Packet) R3() bool {
	return dp.BTN&0b10000000 == 1
}
func (dp DS4_Packet) PS() bool {
	return dp.TOUCHPS&0b1 == 1
}
func (dp DS4_Packet) TOUCH() bool {
	return dp.TOUCHPS&0b10 == 1
}

func (ds4 *DS4) Close() {
	unix.Close(ds4.Ctrl)
	unix.Close(ds4.Data)
}

func writePacket(fd int, pkt *HID_OUTPUT_RESPONSE_PACKET, c chan error) {
	pkt_b := (*[unsafe.Sizeof(*pkt)]byte)(unsafe.Pointer(pkt))[:]
	mycrc := (*uint32)(unsafe.Pointer(&(pkt_b[HID_OUTPUT_RESPONSE_SIZE-4])))
	*mycrc = crc32.ChecksumIEEE(pkt_b[:HID_OUTPUT_RESPONSE_SIZE-4])
	if n, err := unix.Write(fd, pkt_b[:HID_OUTPUT_RESPONSE_SIZE]); err != nil {
		c <- err
	} else if n != HID_OUTPUT_RESPONSE_SIZE {
		c <- fmt.Errorf("Write packet size to DS4 error(should be %d but returns %d).\n", HID_OUTPUT_RESPONSE_SIZE, n)
	}
}

func (ds4 *DS4) SetLED(c color.Color) {
	r, g, b, _ := c.RGBA()
	ds4.SetLEDRGB(int(r), int(g), int(b))
}

func (ds4 *DS4) SetReportType(p int) {
	pkt := initPacket()
	pkt.Protocol = byte(p)
	writePacket(ds4.Ctrl, pkt, ds4.Event)
}

func (ds4 *DS4) SetLEDRGB(r, g, b int) {
	pkt := initPacket()
	pkt.FEATURE = FEATURE_LED
	pkt.LED[0] = byte(r)
	pkt.LED[1] = byte(g)
	pkt.LED[2] = byte(b)
	writePacket(ds4.Ctrl, pkt, ds4.Event)
}

func (ds4 *DS4) SetRumble(powerStrong, powerWeak int) {
	pkt := initPacket()
	pkt.FEATURE = FEATURE_RUMBLE
	pkt.RumbleStrong = byte(powerStrong)
	pkt.RumbleWeak = byte(powerWeak)
	writePacket(ds4.Ctrl, pkt, ds4.Event)
}

func (ds4 *DS4) SetLEDDelay(on, off int) {
	pkt := initPacket()
	pkt.FEATURE = FEATURE_BLINK
	pkt.LEDDelay[0] = byte(on)
	pkt.LEDDelay[1] = byte(off)
	writePacket(ds4.Ctrl, pkt, ds4.Event)
}

func initPacket() *HID_OUTPUT_RESPONSE_PACKET {
	ret := HID_OUTPUT_RESPONSE_PACKET{
		Type:     HID_TRANSACTION_SET | HID_REPORT_OUT,
		Protocol: 0x11,
		HIDCRC:   DTYPE_CRC | DTYPE_HID | DTYPE_INTERVAL_62MS, /*We do not need so fast send*/
	}
	return &ret
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

func NewDS4(addrStr string) (*DS4, error) {
	var ret DS4
	decoded, err := hex.DecodeString(strings.ReplaceAll(addrStr, ":", ""))
	if err != nil {
		return nil, fmt.Errorf("Input bluetooth address %s is invalid.", addrStr)
	}
	var addr [6]uint8
	copy(addr[:], decoded)
	ctl_fd, err := newL2Conn(addr, 0x11)
	if err != nil {
		return nil, err
	}
	data_fd, err := newL2Conn(addr, 0x13)
	if err != nil {
		return nil, err
	}

	ret.IsConn = true
	ret.Ctrl = ctl_fd
	ret.Data = data_fd
	ret.Notify = make(map[DS4NotifyType]chan int)
	// Send 0x11 report
	ret.SetReportType(0x11)

	go func() {
		status_b := (*[unsafe.Sizeof(ret.Status)]byte)(unsafe.Pointer(&ret.Status))[:]
		prev_b := (*[unsafe.Sizeof(ret.PrevStatus)]byte)(unsafe.Pointer(&ret.PrevStatus))[:]
		checkTimer := time.Tick(100 * time.Millisecond)
		isFirst := true
		for {
			select {
			case <-checkTimer:
				ret.CheckNotify()
				copy(prev_b, status_b)
			default:
				n, err := unix.Read(ret.Data, status_b)
				ret.Status.PacketSize = n
				if isFirst {
					copy(prev_b, status_b)
					isFirst = false
				}
				if err != nil {
					ret.Event <- err
					return
				}
				if n == 0 {
					ret.Event <- fmt.Errorf("Data channel Disconnect")
					ret.IsConn = false
					return
				}
				if n == -1 {
					continue
				}
			}
		}
	}()

	return &ret, nil
}
