package main

import (
	"fmt"
	"github.com/eskpil/rc1/internal"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	if 2 > len(os.Args) {
		fmt.Printf("missing arguments: expected operation\n")
		return
	}

	operation := os.Args[1]

	devices, err := internal.ListDevices()
	if err != nil {
		log.Printf("could not list devices: %v\n", err)
		return
	}

	connection, err := internal.DbusConnect()
	if err != nil {
		log.Printf("could not connect with system bus: %v", err)
		return
	}

	adjustmentValue := 20.0
	adjustment := internal.AdjustmentIncrement

	switch operation {
	case "increase":
		{
			adjustment = internal.AdjustmentIncrement

			if 3 > len(os.Args) {
				fmt.Printf("expected percentage\n")
				return
			}

			value, err := strconv.ParseFloat(os.Args[2], 64)
			if err != nil {
				fmt.Printf("could not parse argument as float: %v\n", err)
			}

			adjustmentValue = value
		}
	case "decrease":
		{
			adjustment = internal.AdjustmentDecrement

			if 3 > len(os.Args) {
				fmt.Printf("expected percentage\n")
				return
			}

			value, err := strconv.ParseFloat(os.Args[2], 64)
			if err != nil {
				fmt.Printf("could not parse argument as float: %v\n", err)
			}

			adjustmentValue = value
		}
	case "max":
		{
			adjustment = internal.AdjustmentMax
		}
	case "min":
		{
			adjustment = internal.AdjustmentMin
		}
	case "display":
		{
			maxBrightness, err := devices[0].MaxBrightness()
			if err != nil {
				fmt.Printf("could not get max brightness: %v", err)
				return
			}

			actualBrightness, err := devices[0].ActualBrightness()
			if err != nil {
				fmt.Printf("could not get actual brightness: %v", err)
				return
			}

			percentage := actualBrightness / maxBrightness * 100.0
			fmt.Printf("%d", uint32(math.Abs(percentage)))
			fmt.Println("%")
		}
	default:
		{
			fmt.Printf("invalid arguments: unknown operation: %s\n", operation)
			return
		}
	}

	for _, device := range devices {
		if err := device.Adjust(connection, adjustment, adjustmentValue); err != nil {
			fmt.Printf("could not adjust brightness: %v", err)
		}
	}
}
