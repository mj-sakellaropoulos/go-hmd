package wmr

import (
	"fmt"
	"time"

	"github.com/sstallion/go-hid"
)

/* Ref: monado/src/xrt/drivers/wmr/wmr_common.h */
const (
	MS_HOLOLENS_MANUFACTURER_STRING = "Microsoft"
	MS_HOLOLENS_PRODUCT_STRING      = "HoloLens Sensors"

	MICROSOFT_VID                       = 0x045e
	HOLOLENS_SENSORS_PID                = 0x0659
	WMR_CONTROLLER_PID                  = 0x065b
	WMR_CONTROLLER_LEFT_PRODUCT_STRING  = "Motion controller - Left"
	WMR_CONTROLLER_RIGHT_PRODUCT_STRING = "Motion controller - Right"

	HP_VID                   = 0x03f0
	VR1000_PID               = 0x0367
	REVERB_G1_PID            = 0x0c6a
	REVERB_G2_PID            = 0x0580
	REVERB_G2_CONTROLLER_PID = 0x066a /* On 0x045e Microsoft VID */

	LENOVO_VID   = 0x17ef
	EXPLORER_PID = 0xb801

	DELL_VID  = 0x413c
	VISOR_PID = 0xb0d5

	SAMSUNG_VID            = 0x04e8
	ODYSSEY_PID            = 0x7310
	ODYSSEY_PLUS_PID       = 0x7312
	ODYSSEY_CONTROLLER_PID = 0x065d

	QUANTA_VID              = 0x0408 /* Medion? */
	MEDION_ERAZER_X1000_PID = 0xb5d5
)

const (
	WMR_HEADSET_GENERIC = iota
	WMR_HEADSET_HP_VR1000
	WMR_HEADSET_REVERB_G1
	WMR_HEADSET_REVERB_G2
	WMR_HEADSET_SAMSUNG_XE700X3AI
	WMR_HEADSET_SAMSUNG_800ZAA
	WMR_HEADSET_LENOVO_EXPLORER
	WMR_HEADSET_MEDION_ERAZER_X1000
	WMR_HEADSET_DELL_VISOR
)

type HMD struct {
	headset_type int

	hololens_sensors_device *hid.Device
	companion_device        *hid.Device
	//controller_left_device  *hid.Device
	//controller_right_device *hid.Device
}

func CreateHMD_G2() (*HMD, error) {
	// Open device.
	companion_device, err := hid.OpenFirst(HP_VID, REVERB_G2_PID)
	if err != nil {
		return nil, err
	}

	// Open device.
	sensors_device, err := hid.OpenFirst(MICROSOFT_VID, HOLOLENS_SENSORS_PID)
	if err != nil {
		return nil, err
	}

	hmd := HMD{
		headset_type:            WMR_HEADSET_REVERB_G2,
		hololens_sensors_device: sensors_device,
		companion_device:        companion_device,
	}

	return &hmd, nil
}

func CloseHMDDevices(hmd *HMD) error {
	err := hmd.companion_device.Close()
	if err != nil {
		return err
	}
	err = hmd.hololens_sensors_device.Close()
	if err != nil {
		return err
	}
	return nil
}

/* Ref: monado/src/xrt/drivers/wmr/wmr_hmd.c */
func ActivateReverb(HMD *HMD) {
	//todo: handle errors.
	fmt.Println("Activating HP Reverb G1/G2 HMD...")

	time.Sleep(300 * time.Millisecond)

	for i := 0; i < 4; i++ {
		cmd := make([]byte, 64)
		cmd[0] = 0x50
		cmd[1] = 0x01
		HMD.companion_device.SendFeatureReport(cmd)

		data := make([]byte, 64)
		data[0] = 0x50
		HMD.companion_device.GetFeatureReport(data) //HID_GET(wh, hid, data, "loop")

		time.Sleep(10 * time.Millisecond) //os_nanosleep(U_TIME_1MS_IN_NS * 10) // Sleep 10ms
	}

	data := make([]byte, 64)
	data[0] = 0x09
	HMD.companion_device.GetFeatureReport(data) //HID_GET(wh, hid, data, "data_1")

	data[0] = 0x08
	HMD.companion_device.GetFeatureReport(data) //HID_GET(wh, hid, data, "data_2")

	data[0] = 0x06
	HMD.companion_device.GetFeatureReport(data) //HID_GET(wh, hid, data, "data_3")

	fmt.Println("Sent activation report.")

	ScreenEnableReverb(HMD.companion_device, true)

	fmt.Println("Sleep to wait for display to enumerate.")

	time.Sleep(5 * time.Second)
}

/* Ref: monado/src/xrt/drivers/wmr/wmr_hmd.c */
func ScreenEnableReverb(companion_device *hid.Device, state bool) error {
	// cmd : { 0x04 , 0x00 }
	// cmd[0] feature_number
	// cmd[1] status (0x00 off, 0x01 on)

	cmd := []byte{0x04, 0x01}

	if !state {
		cmd = []byte{0x04, 0x00}
	}

	// Enable Display
	_, err := companion_device.SendFeatureReport(cmd)
	if err != nil {
		return err
	}

	return nil
}
