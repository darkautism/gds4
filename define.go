package go-ds4

const (
	HID_TRANSACTION_GET  byte = 0x40
	HID_TRANSACTION_SET  byte = 0x50
	HID_TRANSACTION_DATA byte = 0xA0
	HID_REPORT_IN        byte = 0x01
	HID_REPORT_OUT       byte = 0x02
	HID_REPORT_FEATURE   byte = 0x03
)

const (
	FEATURE_RUMBLE = 0x01
	FEATURE_LED    = 0x02
	FEATURE_BLINK  = 0x04
)

const (
	// low 6 bit
	DTYPE_INTERVAL_1MS     = 0x00
	DTYPE_INTERVAL_1MS2    = 0x01
	DTYPE_INTERVAL_2MS     = 0x02
	DTYPE_INTERVAL_62MS    = 0x3E
	DTYPE_INTERVAL_DISABLE = 0x3F
	DTYPE_CRC              = 0x40
	DTYPE_HID              = 0x80
)

const HID_OUTPUT_RESPONSE_SIZE = 79

type HID_OUTPUT_RESPONSE_PACKET struct {
	Type         byte
	Protocol     byte
	HIDCRC       byte
	unknow       byte
	FEATURE      byte
	unknow2      [2]byte
	RumbleWeak   byte
	RumbleStrong byte
	LED          [3]byte
	LEDDelay     [2]byte
	unknow3      [65]byte
}

type DS4 struct {
	Data       int
	Ctrl       int
	Status     DS4_Packet
	PrevStatus DS4_Packet
	IsConn     bool
	Notify     map[DS4NotifyType]chan int
	Event      chan error
}

type DS4_Packet struct {
	Type     byte
	Protocol byte
	unknow   byte
	ReportID byte
	L_X      uint8
	L_Y      uint8
	R_X      uint8
	R_Y      uint8
	// DPAD			0b1111
	// SQUARE_PAD	0b10000
	// X_PAD		0b100000
	// O_PAD		0b1000000
	// TRIANGLE_PAD	0b10000000
	PAD uint8
	// L1		0b0001
	// R1		0b0010
	// L2		0b0100
	// R2		0b1000
	// SHARE	0b10000
	// OPTION	0b100000
	// L3		0b1000000
	// R3		0b10000000
	BTN uint8

	// PS		0b0001
	// TOUCH	0b0010
	TOUCHPS         uint8
	L2Analogy       uint8
	R2Analogy       uint8
	Timestamp       [2]byte
	Battery         uint8
	AngularVelocity [3]uint16 /* X,Y,Z */
	Acceleration    [3]uint16 /* X,Y,Z */
	unknow99        [51]byte
	PacketSize      int
}