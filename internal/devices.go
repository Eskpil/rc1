package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Adjustment int

const (
	AdjustmentIncrement Adjustment = iota
	AdjustmentDecrement
	AdjustmentMin
	AdjustmentMax
)

type Device string

func ListDevices() ([]Device, error) {
	var devices []Device

	err := filepath.Walk("/sys/class/backlight", func(path string, info os.FileInfo, err error) error {
		if path == "/sys/class/backlight" {
			return nil
		}

		parts := strings.Split(path, "/")
		device := parts[len(parts)-1]

		devices = append(devices, Device(device))

		return nil
	})

	return devices, err
}

func (d Device) MaxBrightness() (float64, error) {
	contents, err := os.ReadFile(fmt.Sprintf("/sys/class/backlight/%s/max_brightness", d))
	if err != nil {
		return 0, err
	}

	brightness, err := strconv.ParseFloat(string(contents[0:len(contents)-1]), 10)
	if err != nil {
		return 0, err
	}

	return brightness, nil
}

func (d Device) ActualBrightness() (float64, error) {
	contents, err := os.ReadFile(fmt.Sprintf("/sys/class/backlight/%s/actual_brightness", d))
	if err != nil {
		return 0, err
	}

	brightness, err := strconv.ParseFloat(string(contents[0:len(contents)-1]), 10)
	if err != nil {
		return 0, err
	}

	return brightness, nil
}

func (d Device) Adjust(connection *DbusConnection, adjustment Adjustment, adjustmentValue float64) error {
	var value float64 = 0

	actualBrightness, err := d.ActualBrightness()
	if err != nil {
		return err
	}

	switch adjustment {
	case AdjustmentDecrement:
		{
			value = actualBrightness - (adjustmentValue/100)*actualBrightness
		}
	case AdjustmentIncrement:
		{
			value = actualBrightness + (adjustmentValue/100)*actualBrightness
		}
	case AdjustmentMax:
		{
			maxBrightness, err := d.MaxBrightness()
			if err != nil {
				return err
			}

			value = maxBrightness
		}
	case AdjustmentMin:
		{
			value = 0
		}
	}

	session, err := connection.GetActiveSession()
	if err != nil {
		return err
	}

	call := session.Call("org.freedesktop.login1.Session.SetBrightness", 0, "backlight", d, uint32(value))
	if call.Err != nil {
		return err
	}

	return nil
}
