package main

import (
	"fmt"
	"github.com/caseymrm/menuet"
	"github.com/karalabe/hid"
	"log"
	"time"
)

var vendorId uint16 = 0x04d8
var productId uint16 = 0xf372
var currentIcon string = "light-normal.pdf"
var lastCommand []byte
var reSendCommand bool = false

func getLuxaforDevice() *hid.Device {
	for _, info := range hid.Enumerate(vendorId, productId) {
		if info.VendorID == vendorId && info.ProductID == productId {
			dev, err := info.Open()

			if err != nil {
				log.Print(err)
			} else {
				log.Printf("Device connected")
				return dev
			}
		}
	}
	log.Printf("Device not connected")
	return nil
}

func runLuxaforCommand(command []byte) {
	dev := getLuxaforDevice()

	if dev != nil {
		menuet.App().SetMenuState(&menuet.MenuState{
			Image: currentIcon,
		})
		dev.Write(command)
		lastCommand = command
		dev.Close()
	} else {
		menuet.App().SetMenuState(&menuet.MenuState{
			Image: "light-disabled.pdf",
		})
	}
}

func setRGB(red byte, green byte, blue byte) {
	runLuxaforCommand([]byte{1, 255, red, green, blue, 0, 0})
}

func fadeRGB(red byte, green byte, blue byte) {
	runLuxaforCommand([]byte{2, 255, red, green, blue, 50, 0})
	clearMasterColor()
}

func setPattern(pattern byte) {
	runLuxaforCommand([]byte{6, pattern, 255, 0, 0, 0, 0, 0})
	clearMasterColor()
}

func setMasterColor(color string) {
	imageName := fmt.Sprintf("%s.png", color)
	menuet.App().SetMenuState(&menuet.MenuState{
		Image: imageName,
	})
	currentIcon = imageName
}

func clearMasterColor() {
	menuet.App().SetMenuState(&menuet.MenuState{
		Image: "light-normal.pdf",
	})
	currentIcon = "light-normal.pdf"
}

func fadeMenu() []menuet.MenuItem {
	return []menuet.MenuItem{
		menuet.MenuItem{
			Text: "Green",
			Clicked: func() {
				fadeRGB(0, 255, 0)
			},
		},
		menuet.MenuItem{
			Text: "Red",
			Clicked: func() {
				fadeRGB(255, 0, 0)
			},
		},
		menuet.MenuItem{
			Text: "Blue",
			Clicked: func() {
				fadeRGB(0, 0, 255)
			},
		},
	}
}

func patternsMenu() []menuet.MenuItem {
	return []menuet.MenuItem{
		menuet.MenuItem{
			Text: "1",
			Clicked: func() {
				setPattern(1)
			},
		},
		menuet.MenuItem{
			Text: "2",
			Clicked: func() {
				setPattern(2)
			},
		},
		menuet.MenuItem{
			Text: "3",
			Clicked: func() {
				setPattern(3)
			},
		},
		menuet.MenuItem{
			Text: "4",
			Clicked: func() {
				setPattern(4)
			},
		},
		menuet.MenuItem{
			Text: "5",
			Clicked: func() {
				setPattern(5)
			},
		},
		menuet.MenuItem{
			Text: "6",
			Clicked: func() {
				setPattern(6)
			},
		},
	}
}

func menuItems() []menuet.MenuItem {
	return []menuet.MenuItem{
		menuet.MenuItem{
			Text: "Green",
			Clicked: func() {
				setRGB(0, 255, 0)
				setMasterColor("green")
			},
		},
		menuet.MenuItem{
			Text: "Yellow",
			Clicked: func() {
				setRGB(255, 255, 0)
				setMasterColor("yellow")
			},
		},
		menuet.MenuItem{
			Text: "Red",
			Clicked: func() {
				setRGB(255, 0, 0)
				setMasterColor("red")
			},
		},
		menuet.MenuItem{
			Text: "Blue",
			Clicked: func() {
				setRGB(0, 0, 255)
				setMasterColor("blue")
			},
		},
		menuet.MenuItem{
			Text:     "Fade",
			Children: fadeMenu,
		},
		menuet.MenuItem{
			Text:     "Patterns",
			Children: patternsMenu,
		},
		menuet.MenuItem{
			Text: "Off",
			Clicked: func() {
				setRGB(0, 0, 0)
				clearMasterColor()
			},
		},
	}
}

func updateDeviceStatus() {
	dev := getLuxaforDevice()

	if dev != nil {
		log.Print(currentIcon)
		dev.Close()
		if reSendCommand && len(lastCommand) > 0 {
			runLuxaforCommand(lastCommand)
		}
		reSendCommand = false
	} else {
		menuet.App().SetMenuState(&menuet.MenuState{
			Image: "light-disabled.pdf",
		})
		reSendCommand = true
	}
}

func main() {
	app := menuet.App()

	app.SetMenuState(&menuet.MenuState{
		Image: "light-normal.pdf",
	})

	app.Children = menuItems
	updateDeviceStatus()

	app.Name = "Luxafor macOS"
	app.Label = "luxafor-macos.colde.github.com"

	// Setup device update ticket
	ticker := time.NewTicker(20 * time.Second)
	tickerDone := make(chan bool)

	go func() {
		for {
			select {
			case <-tickerDone:
				return
			case _ = <-ticker.C:
				log.Printf("Updating device status")
				updateDeviceStatus()
			}
		}
	}()

	app.RunApplication()
	ticker.Stop()
	tickerDone <- true
}
