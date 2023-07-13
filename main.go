package main

/*
#cgo CFLAGS: -I ./include
#cgo LDFLAGS: -L ./lib -lvatek_core -lusb-1.0

#include <stdint.h>

int GetVatekSDKVersion();
char* NewVatekContext();
int FreeVatekContext(char* p);
int VatekDeviceOpen(char* p);
int GetVatekDeviceChipInfo(char* p, int* status, uint32_t* fwVer, uint32_t* chipId, uint32_t* service, uint32_t* in, uint32_t* out, uint32_t* peripheral);
*/
import "C"
import "fmt"

type VatekContext *C.char

type VatekChipInfo struct {
	Status        ChipStatus
	FwVer         uint32
	ChipModule    uint32
	Service       uint32
	InputSupport  uint32
	OutputSupport uint32
	Peripheral    uint32
}

func GetVatekSDKVersion() int {
	return int(C.GetVatekSDKVersion())
}

func GetVatekDeviceChipInfo(ctx VatekContext) VatekChipInfo {
	var status C.int
	var fwVer, chipId, service, in, out, peripheral C.uint
	C.GetVatekDeviceChipInfo(ctx, &status, &fwVer, &chipId, &service, &in, &out, &peripheral)
	return VatekChipInfo{
		Status:        ChipStatus(status),
		FwVer:         uint32(fwVer),
		ChipModule:    uint32(chipId),
		Service:       uint32(service),
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
	fmt.Println("-------------------------------------")
	fmt.Println("	   Chip Information")
	fmt.Println("-------------------------------------")
	fmt.Printf("%-20s : %s\n", "Chip Status", chip.Status.String())
	fmt.Printf("%-20s : %08x\n", "FW Version", chip.FwVer)
	fmt.Printf("%-20s : %08x\n", "Chip  ID", chip.ChipModule)
	fmt.Printf("%-20s : %08x\n", "Service", chip.Service)
	fmt.Printf("%-20s : %08x\n", "Input", chip.InputSupport)
	fmt.Printf("%-20s : %08x\n", "Output", chip.OutputSupport)
	fmt.Printf("%-20s : %08x\n", "Peripheral", chip.Peripheral)
	fmt.Printf("%-20s : %d\n", "SDK Version", GetVatekSDKVersion())
}
