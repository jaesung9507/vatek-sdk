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
	case -4:
		return "badstatus"
	case -3:
		return "fail_hw"
	case -2:
		return "fail_service"
	case -1:
		return "fail_loader"
	case 1:
		return "waitcmd"
	case 2:
		return "running"
	}
	return "unknown"
}
