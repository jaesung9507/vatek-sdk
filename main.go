package main

/*
#cgo CFLAGS: -I ./include
#cgo LDFLAGS: -L ./lib -lvatek_core -lusb-1.0

#include <stdint.h>

int GetVatekSDKVersion();
char* NewVatekContext();
int FreeVatekContext(char* p);
int VatekDeviceOpen(char* p);
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

type VatekContext *C.char

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

func GetVatekDeviceChipInfo(ctx VatekContext) VatekChipInfo {
	var status, chipId C.int
	var fwVer, service, in, out, peripheral C.uint
	C.GetVatekDeviceChipInfo(ctx, &status, &fwVer, &chipId, &service, &in, &out, &peripheral)
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

func GetVatekUsbStreamStatus(ctx VatekContext) (UsbStreamStatus, TransformInfo) {
	var status C.int = -1
	var cur, data, mode C.uint
	C.GetVatekUsbStreamStatus(ctx, &status, &cur, &data, &mode)
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
	ctx := VatekContext(C.NewVatekContext())
	defer C.FreeVatekContext(ctx)
	errNum := C.VatekDeviceOpen(ctx)
	if errNum < 0 {
		fmt.Printf("failed to device open: %d\n", errNum)
		return
	}
	chip := GetVatekDeviceChipInfo(ctx)
	fmt.Print(chip.String())

	errNum = C.VatekUsbStreamOpen(ctx)
	if errNum < 0 {
		fmt.Printf("failed to usb stream open: %d\n", errNum)
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

	C.VatekUsbStreamStart(ctx)
	if errNum < 0 {
		fmt.Printf("failed to usb stream open: %d\n", errNum)
		return
	}
	errCnt := 0
	tick := time.Now()
	for {
		status, info := GetVatekUsbStreamStatus(ctx)
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
