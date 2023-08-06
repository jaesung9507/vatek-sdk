package vatek

/*
#cgo CFLAGS: -I /usr/local/include/vatek
#cgo LDFLAGS: -lvatek_core -lusb-1.0

#include <stdint.h>

int GetVatekSDKVersion();
char* NewVatekContext(int modulatorType, uint32_t freqkhz);
int FreeVatekContext(char* p);
int VatekUsbDeviceOpen(char* p);
int GetVatekDeviceChipInfo(char* p, int* status, uint32_t* fwVer, int* chipId, uint32_t* service, uint32_t* in, uint32_t* out, uint32_t* peripheral);
int VatekUsbStreamOpen(char *p);
int SetVatekCallbackParam(char* p, void* param);
int VatekUsbStreamStart(char* p);
int GetVatekUsbStreamStatus(char* p, int* status, uint32_t* cur, uint32_t* data, uint32_t* mode);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type VatekContext struct {
	p           *C.char
	getTsStream func() []byte
}

type VatekChipInfo struct {
	Status        ChipStatus
	FwVer         uint32
	ChipModule    ChipID
	Service       ServiceMode
	InputSupport  uint32
	OutputSupport uint32
	Peripheral    uint32
}

type BroadcastInfo struct {
	CurBitrate  uint32
	DataBitrate uint32
}

type TransformInfo struct {
	Info BroadcastInfo
	Mode TransformMode
}

func NewVatekContext(modulatorType ModulatorType, freqkhz uint32) VatekContext {
	return VatekContext{p: C.NewVatekContext(C.int(modulatorType), C.uint(freqkhz))}
}

func (ctx *VatekContext) Close() {
	C.FreeVatekContext(ctx.p)
}

func (ctx *VatekContext) UsbDeviceOpen() error {
	errNum := C.VatekUsbDeviceOpen(ctx.p)
	if errNum < 0 {
		return VatekError(errNum)
	}
	return nil
}

func (ctx *VatekContext) UsbStreamOpen() error {
	errNum := C.VatekUsbStreamOpen(ctx.p)
	if errNum < 0 {
		return VatekError(errNum)
	}
	return nil
}

func (ctx *VatekContext) UsbStreamStart(getTsFrame func() []byte) error {
	if ctx.getTsStream == nil {
		C.SetVatekCallbackParam(ctx.p, unsafe.Pointer(ctx))
	}
	ctx.getTsStream = getTsFrame
	errNum := C.VatekUsbStreamStart(ctx.p)
	if errNum < 0 {
		return VatekError(errNum)
	}

	return nil
}

func (c *VatekChipInfo) String() string {
	str := "-------------------------------------\n"
	str += "	   Chip Information\n"
	str += "-------------------------------------\n"
	str += fmt.Sprintf("%-20s : %s\n", "Chip Status", c.Status.String())
	str += fmt.Sprintf("%-20s : %08x\n", "FW Version", c.FwVer)
	str += fmt.Sprintf("%-20s : %s\n", "Chip ID", c.ChipModule.String())
	str += fmt.Sprintf("%-20s : %s\n", "Service", c.Service.String())
	str += fmt.Sprintf("%-20s : %08x\n", "Input", c.InputSupport)
	str += fmt.Sprintf("%-20s : %08x\n", "Output", c.OutputSupport)
	str += fmt.Sprintf("%-20s : %08x\n", "Peripheral", c.Peripheral)
	str += fmt.Sprintf("%-20s : %d\n", "SDK Version", GetVatekSDKVersion())
	str += "=====================================\n"
	return str
}

func GetVatekSDKVersion() int {
	return int(C.GetVatekSDKVersion())
}

func (ctx *VatekContext) GetDeviceChipInfo() VatekChipInfo {
	var status, chipId C.int
	var fwVer, service, in, out, peripheral C.uint
	C.GetVatekDeviceChipInfo(ctx.p, &status, &fwVer, &chipId, &service, &in, &out, &peripheral)
	return VatekChipInfo{
		Status:        ChipStatus(status),
		FwVer:         uint32(fwVer),
		ChipModule:    ChipID(chipId),
		Service:       ServiceMode(service),
		InputSupport:  uint32(in),
		OutputSupport: uint32(out),
		Peripheral:    uint32(peripheral),
	}
}

func (ctx *VatekContext) GetUsbStreamStatus() (UsbStreamStatus, TransformInfo) {
	var status C.int = -1
	var cur, data, mode C.uint
	C.GetVatekUsbStreamStatus(ctx.p, &status, &cur, &data, &mode)
	tinfo := TransformInfo{
		Info: BroadcastInfo{
			CurBitrate:  uint32(cur),
			DataBitrate: uint32(data),
		},
		Mode: TransformMode(mode),
	}
	return UsbStreamStatus(status), tinfo
}

//export GetTsFrame
func GetTsFrame(param unsafe.Pointer) *C.uchar {
	ctx := (*VatekContext)(param)
	return (*C.uchar)(C.CBytes(ctx.getTsStream()))
}
