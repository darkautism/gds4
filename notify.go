package gds4

/* Because ds4's report fequency is too quickly
 * to cause notify trigger too quickly.
 * This feature is cost too many cpu on my beaglebone
 * and waste a lot of battery.
 */

type DS4NotifyType int

const (
	DS4NotifyTypeLAnalogy     DS4NotifyType = 0
	DS4NotifyTypeRAnalogy     DS4NotifyType = 1
	DS4NotifyTypeDPAD_UP      DS4NotifyType = 2
	DS4NotifyTypeDPAD_DOWN    DS4NotifyType = 3
	DS4NotifyTypeDPAD_LEFT    DS4NotifyType = 4
	DS4NotifyTypeDPAD_RIGHT   DS4NotifyType = 5
	DS4NotifyTypeBTN_SQUARE   DS4NotifyType = 6
	DS4NotifyTypeBTN_X        DS4NotifyType = 7
	DS4NotifyTypeBTN_O        DS4NotifyType = 8
	DS4NotifyTypeBTN_TRIANGLE DS4NotifyType = 9
	DS4NotifyTypeBTN_L1       DS4NotifyType = 11
	DS4NotifyTypeBTN_R1       DS4NotifyType = 12
	DS4NotifyTypeBTN_L2       DS4NotifyType = 13
	DS4NotifyTypeBTN_R2       DS4NotifyType = 14
	DS4NotifyTypeBTN_SHARE    DS4NotifyType = 15
	DS4NotifyTypeBTN_OPTION   DS4NotifyType = 16
	DS4NotifyTypeBTN_L3       DS4NotifyType = 17
	DS4NotifyTypeBTN_R3       DS4NotifyType = 18
	DS4NotifyTypeTOUCHPAD     DS4NotifyType = 19
	DS4NotifyTypePS           DS4NotifyType = 20
	DS4NotifyTypeL2Analogy    DS4NotifyType = 21
	DS4NotifyTypeR2Analogy    DS4NotifyType = 22
)

func (ds4 *DS4) AddNotify(t DS4NotifyType, c chan int) {
	ds4.Notify[t] = c
}

func (ds4 *DS4) DelNotify(t DS4NotifyType) {
	ds4.Notify[t] = nil
	delete(ds4.Notify, t)
}

func (ds4 *DS4) CheckNotify() {
	for t, c := range ds4.Notify {
		switch t {
		case DS4NotifyTypeLAnalogy:
			change_x := int(ds4.PrevStatus.L_X) - int(ds4.Status.L_X)
			change_y := int(ds4.PrevStatus.L_Y) - int(ds4.Status.L_Y)
			if change_x > 2 || change_x < -2 || change_y > 2 || change_y < -2 {
				c <- int(ds4.Status.L_X)<<8 | int(ds4.Status.L_Y)
			}
			break
		case DS4NotifyTypeRAnalogy:
			change_x := int(ds4.PrevStatus.R_X) - int(ds4.Status.R_X)
			change_y := int(ds4.PrevStatus.R_Y) - int(ds4.Status.R_Y)
			if change_x > 2 || change_x < -2 || change_y > 2 || change_y < -2 {
				c <- int(ds4.Status.R_X)<<4 | int(ds4.Status.R_Y)
			}
			break
		case DS4NotifyTypeDPAD_UP:
			if ds4.PrevStatus.Up() != ds4.Status.Up() {
				c <- int(ds4.Status.PAD & 0xF)
			}
			break
		case DS4NotifyTypeDPAD_DOWN:
			if ds4.PrevStatus.Down() != ds4.Status.Down() {
				c <- int(ds4.Status.PAD & 0xF)
			}
			break
		case DS4NotifyTypeDPAD_LEFT:
			if ds4.PrevStatus.Left() != ds4.Status.Left() {
				c <- int(ds4.Status.PAD & 0xF)
			}
			break
		case DS4NotifyTypeDPAD_RIGHT:
			if ds4.PrevStatus.Right() != ds4.Status.Right() {
				c <- int(ds4.Status.PAD & 0xF)
			}
			break
		case DS4NotifyTypeBTN_SQUARE:
			if (ds4.PrevStatus.PAD & 0x10) != (ds4.Status.PAD & 0x10) {
				c <- int(ds4.Status.PAD & 0x10)
			}
			break
		case DS4NotifyTypeBTN_X:
			if (ds4.PrevStatus.PAD & 0x20) != (ds4.Status.PAD & 0x20) {
				c <- int(ds4.Status.PAD & 0x20)
			}
			break
		case DS4NotifyTypeBTN_O:
			if (ds4.PrevStatus.PAD & 0x40) != (ds4.Status.PAD & 0x40) {
				c <- int(ds4.Status.PAD & 0x40)
			}
			break
		case DS4NotifyTypeBTN_TRIANGLE:
			if (ds4.PrevStatus.PAD & 0x80) != (ds4.Status.PAD & 0x80) {
				c <- int(ds4.Status.PAD & 0x80)
			}
			break
		case DS4NotifyTypeBTN_L1:
			if (ds4.PrevStatus.BTN & 0x01) != (ds4.Status.BTN & 0x01) {
				c <- int(ds4.Status.BTN & 0x1)
			}
			break
		case DS4NotifyTypeBTN_R1:
			if (ds4.PrevStatus.BTN & 0x02) != (ds4.Status.BTN & 0x02) {
				c <- int(ds4.Status.BTN & 0x2)
			}
			break
		case DS4NotifyTypeBTN_L2:
			if (ds4.PrevStatus.BTN & 0x04) != (ds4.Status.BTN & 0x04) {
				c <- int(ds4.Status.BTN & 0x4)
			}
			break
		case DS4NotifyTypeBTN_R2:
			if (ds4.PrevStatus.BTN & 0x08) != (ds4.Status.BTN & 0x08) {
				c <- int(ds4.Status.BTN & 0x8)
			}
			break
		case DS4NotifyTypeBTN_SHARE:
			if (ds4.PrevStatus.BTN & 0x10) != (ds4.Status.BTN & 0x10) {
				c <- int(ds4.Status.BTN & 0x10)
			}
			break
		case DS4NotifyTypeBTN_OPTION:
			if (ds4.PrevStatus.BTN & 0x20) != (ds4.Status.BTN & 0x20) {
				c <- int(ds4.Status.BTN & 0x20)
			}
			break
		case DS4NotifyTypeBTN_L3:
			if (ds4.PrevStatus.BTN & 0x40) != (ds4.Status.BTN & 0x40) {
				c <- int(ds4.Status.BTN & 0x40)
			}
			break
		case DS4NotifyTypeBTN_R3:
			if (ds4.PrevStatus.BTN & 0x80) != (ds4.Status.BTN & 0x80) {
				c <- int(ds4.Status.BTN & 0x80)
			}
			break
		case DS4NotifyTypeTOUCHPAD:
			if (ds4.PrevStatus.TOUCHPS & 0x02) != (ds4.Status.TOUCHPS & 0x02) {
				c <- int(ds4.Status.TOUCHPS & 0x02)
			}
			break
		case DS4NotifyTypePS:
			if (ds4.PrevStatus.TOUCHPS & 0x01) != (ds4.Status.TOUCHPS & 0x01) {
				c <- int(ds4.Status.TOUCHPS & 0x01)
			}
			break
		case DS4NotifyTypeL2Analogy:
			if ds4.PrevStatus.L2Analogy != ds4.Status.L2Analogy {
				c <- int(ds4.Status.L2Analogy)
			}
			break
		case DS4NotifyTypeR2Analogy:
			if ds4.PrevStatus.R2Analogy != ds4.Status.R2Analogy {
				c <- int(ds4.Status.R2Analogy)
			}
			break
		}
	}
}
