package main

import (
	"fmt"
	"io"
	"os"
	"time"

	vatek "github.com/jaesung9507/vatek-sdk"
)

func main() {
	ctx := vatek.NewVatekContext(vatek.ModulatorATSC, 473000)
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

	frames, _ := io.ReadAll(f)
	if err = ctx.UsbStreamStart(func() []byte {
		buf := frames[:vatek.ChipStreamSliceLen]
		frames = frames[vatek.ChipStreamSliceLen:]
		frames = append(frames, buf...)
		return buf
	}); err != nil {
		fmt.Printf("failed to usb stream open: %s\n", err.Error())
		return
	}
	errCnt := 0
	tick := time.Now()
	for {
		status, info := ctx.GetUsbStreamStatus()
		if status == vatek.UsbStreamStatusRunning {
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
