package main

/*
#cgo CFLAGS: -I ./include
#cgo LDFLAGS: -L ./lib -lvatek_core -lusb-1.0

char* NewVatekContext();
void FreeVatekContext(char* p);
*/
import "C"

func main() {
	ctx := C.NewVatekContext()
	C.FreeVatekContext(ctx)
}
