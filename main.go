package main

/*
#cgo CFLAGS: -I ./include
#cgo LDFLAGS: -L ./lib -lvatek_core -lusb-1.0

#include <stdint.h>

int GetVatekSDKVersion();
char* NewVatekContext();
int FreeVatekContext(char* p);
int VatekUsbDeviceOpen(char* p);
int GetVatekDeviceChipInfo(char* p, int* status, uint32_t* fwVer, int* chipId, uint32_t* service, uint32_t* in, uint32_t* out, uint32_t* peripheral);
int VatekUsbStreamOpen(char *p);
int VatekUsbStreamStart(char* p);
int GetVatekUsbStreamStatus(char* p, int* status, uint32_t* cur, uint32_t* data, uint32_t* mode);
*/
import "C"
import (
	"fmt"
	"io"
	"os"
	"time"
)

type VatekContext struct {
	p *C.char
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

func NewVatekContext() VatekContext {
	return VatekContext{p: C.NewVatekContext()}
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

func (ctx *VatekContext) UsbStreamStart() error {
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

func main() {
	ctx := NewVatekContext()
	defer ctx.Close()
	err := ctx.UsbDeviceOpen()
	if err != nil {
		fmt.Printf("failed to device open: %s\n", err.Error())
		return
	}
	chip := ctx.GetDeviceChipInfo()
	fmt.Print(chip.String())

	if err = ctx.UsbStreamOpen(); err != nil {
		fmt.Printf("failed to usb stream open: %s\n", err.Error())
		return
	}

	filename := "sample.ts"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	buf, _ := io.ReadAll(f)
	for len(buf) >= 24064 {
		frame := buf[:24064]
		frames = append(frames, frame)
		buf = buf[24064:]
	}

	if err = ctx.UsbStreamStart(); err != nil {
		fmt.Printf("failed to usb stream open: %s\n", err.Error())
		return
	}
	errCnt := 0
	tick := time.Now()
	for {
		status, info := ctx.GetUsbStreamStatus()
		if status == UsbStreamStatusRunning {
			if time.Since(tick) > time.Second {
				tick = time.Now()
				fmt.Printf("Data:[%d]  Current:[%d]\n", info.Info.DataBitrate, info.Info.CurBitrate)
				if info.Info.DataBitrate == 0 || info.Info.CurBitrate == 0 {
					errCnt++
				}
				if errCnt >= 30 {
					break
				}
			}
		} else {
			break
		}
		time.Sleep(time.Millisecond)
	}
}

var index = 0
var frames [][]byte

//export GetTsFrame
func GetTsFrame() *C.uchar {
	p := (*C.uchar)(C.CBytes(frames[index]))
	index++
	if index >= len(frames) {
		index = 0
	}

	return p
}
