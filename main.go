package main

import (
	"fmt"

	"github.com/mj-sakellaropoulos/go-hmd/pkg/wmr"
	"github.com/sstallion/go-hid"
)

// main
func main() {

	if err := hid.Init(); err != nil {
		panic(err)
	}

	HMD, err := wmr.CreateHMD_G2()

	if err != nil {
		fmt.Printf("Failed to create HMD: %s", err)
		panic(err)
	}

	// Activate headset.
	wmr.ActivateReverb(HMD)

	//wmr.ScreenEnableReverb(companion_device, false) // turns off screen

	defer wmr.CloseHMDDevices(HMD)
	defer hid.Exit()
}

func printAllHID() {
	hid.Enumerate(hid.VendorIDAny, hid.ProductIDAny, func(info *hid.DeviceInfo) error {
		fmt.Printf("Found: %d:%d , %s | %s \n", info.VendorID, info.ProductID, info.MfrStr, info.ProductStr)
		return nil
	})
}
