package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/gousb"
)

const (
	vendorT1            = 0x534c
	productT1Bootloader = 0x0000
	productT1Firmware   = 0x0001
	vendorT2            = 0x1209
	productT2Bootloader = 0x53C0
	productT2Firmware   = 0x53C1
	webEpIn             = 0x81
	webEpOut            = 0x01
)

func main() {
	ctx := gousb.NewContext()
	readFeatures(ctx)
	// t2 stops in the second iteration
	// t1 works
	readFeatures(ctx)
}

func readFeatures(ctx *gousb.Context) {

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		if match(desc) {
			return true
		}
		return false
	})

	if err != nil {
		panic(err)
	}
	if len(devs) <= 0 {
		log.Printf("No device")
		return
	}

	device := devs[0]

	iface, done, err := device.DefaultInterface()
	if err != nil {
		panic(err)
	}

	outEndpoint, err := iface.OutEndpoint(webEpOut)
	if err != nil {
		panic(err)
	}

	inEndpoint, err := iface.InEndpoint(webEpIn)
	if err != nil {
		panic(err)
	}

	initialize := &Message{
		Kind: 0,
		Data: []byte{},
	}

	initialize.WriteTo(outEndpoint)
	initialize.ReadFrom(inEndpoint)

	log.Printf(spew.Sdump(initialize))

	done()
	err = device.Close()
	if err != nil {
		panic(err)
	}
}

func match(desc *gousb.DeviceDesc) bool {
	vid := desc.Vendor
	pid := desc.Product
	trezor1 := vid == vendorT1 && (pid == productT1Firmware || pid == productT1Bootloader)
	trezor2 := vid == vendorT2 && (pid == productT2Firmware || pid == productT2Bootloader)
	return trezor1 || trezor2
}
