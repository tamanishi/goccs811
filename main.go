package main

import (
	"fmt"
	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/experimental/devices/ccs811"
	"periph.io/x/periph/host"
)

const (
	ccs811WaitAfterAppStart = 1000 * time.Microsecond // The CCS811 needs a wait after app start

	ccs811RegisterStatus = 0x00

	ccs811RegisterFWAppVersion = 0x24 // 2 bytes
	ccs811RegisterAppStart     = 0xF4 // 0 bytes
)

func main() {
	_, err := host.Init()
	if err != nil {
		panic(err)
	}

	bus, err := i2creg.Open("")
	if err != nil {
		panic(err)
	}

	// dev := i2c.Dev{Bus: bus, Addr: 0x5a}

	// data := i2cRead(&dev, ccs811RegisterFWAppVersion, 2)
	// fmt.Printf("%02x\n", data[0])
	// fmt.Printf("%02x\n", data[1])

	// readDeviceStatus(&dev)

	// i2cWrite(&dev, ccs811RegisterAppStart, []byte{})
	// time.Sleep(ccs811WaitAfterAppStart)

	// readDeviceStatus(&dev)

	dev, err := ccs811.New(bus, &ccs811.DefaultOpts)
	if err != nil {
		panic(err)
	}

	modeWrite := ccs811.MeasurementModeParams{
		MeasurementMode:   ccs811.MeasurementModeConstant1000,
		GenerateInterrupt: false,
		UseThreshold:      false,
	}

	err = dev.SetMeasurementModeRegister(modeWrite)
	if err != nil {
		panic(err)
	}

	modeRead, err := dev.GetMeasurementModeRegister()
	if err != nil {
		panic(err)
	}

	switch modeRead.MeasurementMode {
	case ccs811.MeasurementModeIdle:
		fmt.Println("Idle, low power mode")
	case ccs811.MeasurementModeConstant1000:
		fmt.Println("Constant power mode, IAQ measurement every second")
	case ccs811.MeasurementModePulse:
		fmt.Println("Pulse heating mode IAQ measurement every 10 seconds")
	case ccs811.MeasurementModeLowPower:
		fmt.Println("Low power pulse heating mode IAQ measurement every 60 seconds")
	case ccs811.MeasurementModeConstant250:
		fmt.Println("Constant power mode, sensor measurement every 250ms")
	default:
		fmt.Println("Unknown")
	}

	status, err := dev.ReadStatus()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status: %x\n", status)
	var val = ccs811.SensorValues{}
	err = dev.Sense(&val)
	// err = dev.SensePartial(ccs811.ReadAll, &val)
	if err != nil {
		panic(err)
	}

	baseline, err := dev.GetBaseline()
	if err != nil {
		panic(err)
	}

	fmt.Println(val)

	fmt.Printf("Baseline: %d\n", baseline)
	fmt.Printf("ECO2: %d\n", val.ECO2)
	fmt.Printf("VOC: %d\n", val.VOC)
}

func readDeviceStatus(dev *i2c.Dev) {
	data := i2cRead(dev, ccs811RegisterStatus, 1)
	fmt.Printf("I2C: Device status %08b\n", data[0])
}

func i2cRead(dev *i2c.Dev, register byte, outBytes int) []byte {
	data := make([]byte, outBytes)
	err := dev.Tx([]byte{register}, data)
	if err != nil {
		panic(err)
	}
	return data
}

func i2cWrite(dev *i2c.Dev, register byte, bytesToWrite []byte) {
	data := append([]byte{register}, bytesToWrite...)
	err := dev.Tx(data, nil)
	if err != nil {
		panic(err)
	}
}
