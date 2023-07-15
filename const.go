package main

type ChipStatus int

const (
	ChipStatusBadStatus   = ChipStatus(-4)
	ChipStatusFailHw      = ChipStatus(-3)
	ChipStatusFailService = ChipStatus(-2)
	ChipStatusFailLoader  = ChipStatus(-1)
	ChipStatusUnknown     = ChipStatus(0)
	ChipStatusWaitCmd     = ChipStatus(1)
	ChipStatusRunning     = ChipStatus(2)
)

func (s ChipStatus) String() string {
	switch s {
	case ChipStatusBadStatus:
		return "badstatus"
	case ChipStatusFailHw:
		return "fail_hw"
	case ChipStatusFailService:
		return "fail_service"
	case ChipStatusFailLoader:
		return "fail_loader"
	case ChipStatusWaitCmd:
		return "waitcmd"
	case ChipStatusRunning:
		return "running"
	}
	return "unknown"
}

type ChipID int

const (
	ChipIdNoDevice = ChipID(-1)
	ChipIdA1       = ChipID(0x00010100)
	ChipIdB1       = ChipID(0x00020100)
	ChipIdB2       = ChipID(0x00020200)
	ChipIdB2Plus   = ChipID(0x00020201)
	ChipIdA3       = ChipID(0x00010300)
	ChipIdB3Lite   = ChipID(0x00020300)
	ChipIdB3Plus   = ChipID(0x00020301)
	ChipIdE1       = ChipID(0x00040300)
	ChipIdUnknown  = ChipID(0x00ffff00)
)

func (c ChipID) String() string {
	switch c {
	case ChipIdNoDevice:
		return "ic_nodevice"
	case ChipIdA1:
		return "a1"
	case ChipIdB1:
		return "b1"
	case ChipIdB2:
		return "b2"
	case ChipIdA3:
		return "a3"
	case ChipIdB2Plus:
		return "b2_plus"
	case ChipIdB3Lite:
		return "b3_lite"
	case ChipIdB3Plus:
		return "b3_plus"
	}
	return "ic_unknown"
}

type ServiceMode uint32

const (
	SeviceUnknown    = ServiceMode(0)
	ServcieRescue    = ServiceMode(0xFF000001)
	ServcieBroadcast = ServiceMode(0xF8000001)
	ServcieTransform = ServiceMode(0xF8000002)
)

func (m ServiceMode) String() string {
	switch m {
	case ServcieRescue:
		return "rescue"
	case ServcieBroadcast:
		return "broadcast"
	case ServcieTransform:
		return "transform"
	}
	return "unknown"
}

type UsbStreamStatus int

const (
	UsbStreamErrUnknown     = UsbStreamStatus(-1)
	UsbStreamStatusIdle     = UsbStreamStatus(0)
	UsbStreamStatusRunning  = UsbStreamStatus(1)
	UsbStreamStatusMoredata = UsbStreamStatus(2)
	UsbStreamStatusStopping = UsbStreamStatus(3)
	UsbStreamStatusStop     = UsbStreamStatus(4)
)

type TransformMode uint32

const (
	TransformNULL      = TransformMode(0)
	TransformEnum      = TransformMode(0x00000003)
	TransformCapture   = TransformMode(0x00000004)
	TransformBroadcast = TransformMode(0x00000005)
)

func (m TransformMode) String() string {
	switch m {
	case TransformEnum:
		return "enum"
	case TransformCapture:
		return "capture"
	case TransformBroadcast:
		return "broadcast"
	}
	return ""
}
