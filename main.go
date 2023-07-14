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
*/
import "C"
import "fmt"

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

func main() {
	ctx := VatekContext(C.NewVatekContext())
	defer C.FreeVatekContext(ctx)
	err := C.VatekDeviceOpen(ctx)
	if err < 0 {
		fmt.Printf("failed to device open: %d\n", err)
		return
	}
	chip := GetVatekDeviceChipInfo(ctx)
	fmt.Print(chip.String())
}
