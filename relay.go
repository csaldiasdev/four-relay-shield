package main

import (
	"sync"
	"time"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

var once sync.Once

type gpiodProperties struct {
	relay4line *gpiod.Line
	relay3line *gpiod.Line
	relay2line *gpiod.Line
	relay1line *gpiod.Line
	chip       *gpiod.Chip
}

// ChangeRelayStatus struct
type ChangeRelayStatus struct {
	relayNumber   int
	relayNewValue bool
	changedAt     time.Time
}

// RelayController struct
type RelayController struct {
	gpiodProperties
	changeStatusChannel chan (ChangeRelayStatus)
}

var controllerInstance *RelayController

// SetRelay1Status allows set relay #1 status
func (gc RelayController) SetRelay1Status(status bool) {
	if gc.relay1line.SetValue(boolToInt(status)) == nil {
		gc.changeStatusChannel <- ChangeRelayStatus{
			relayNumber:   1,
			relayNewValue: status,
			changedAt:     time.Now(),
		}
	}
}

// SetRelay2Status allows set relay #2 status
func (gc RelayController) SetRelay2Status(status bool) {
	if gc.relay2line.SetValue(boolToInt(status)) == nil {
		gc.changeStatusChannel <- ChangeRelayStatus{
			relayNumber:   2,
			relayNewValue: status,
			changedAt:     time.Now(),
		}
	}
}

// SetRelay3Status allows set relay #3 status
func (gc RelayController) SetRelay3Status(status bool) {
	if gc.relay3line.SetValue(boolToInt(status)) == nil {
		gc.changeStatusChannel <- ChangeRelayStatus{
			relayNumber:   3,
			relayNewValue: status,
			changedAt:     time.Now(),
		}
	}
}

// SetRelay4Status allows set relay #4 status
func (gc RelayController) SetRelay4Status(status bool) {
	if gc.relay4line.SetValue(boolToInt(status)) == nil {
		gc.changeStatusChannel <- ChangeRelayStatus{
			relayNumber:   4,
			relayNewValue: status,
			changedAt:     time.Now(),
		}
	}
}

// GetAllRelaysStatus returns all relays status
func (gc RelayController) GetAllRelaysStatus() (bool, bool, bool, bool) {
	relay1value, _ := gc.relay1line.Value()
	relay2value, _ := gc.relay2line.Value()
	relay3value, _ := gc.relay3line.Value()
	relay4value, _ := gc.relay4line.Value()

	return intToBool(relay1value), intToBool(relay2value), intToBool(relay3value), intToBool(relay4value)
}

// Close releases all resources
func (gc RelayController) Close() {
	close(gc.changeStatusChannel)
	gc.relay1line.Close()
	gc.relay2line.Close()
	gc.relay3line.Close()
	gc.relay4line.Close()
	gc.chip.Close()
}

// InitController initialize the controller
func InitController() error {
	if controllerInstance == nil {
		once.Do(func() {

			cc := gpiod.Chips()

			c, _ := gpiod.NewChip(cc[0])

			relay1line, _ := c.RequestLine(rpi.J8p15, gpiod.AsOutput(0)) // RELAY 1
			relay2line, _ := c.RequestLine(rpi.J8p13, gpiod.AsOutput(0)) // RELAY 2
			relay3line, _ := c.RequestLine(rpi.J8p11, gpiod.AsOutput(0)) // RELAY 3
			relay4line, _ := c.RequestLine(rpi.J8p7, gpiod.AsOutput(0))  // RELAY 4

			controllerInstance = &RelayController{
				changeStatusChannel: make(chan ChangeRelayStatus),
				gpiodProperties: gpiodProperties{
					chip:       c,
					relay1line: relay1line,
					relay2line: relay2line,
					relay3line: relay3line,
					relay4line: relay4line,
				},
			}
		})
	}

	return nil
}

// GetControllerInstance gets relay controller instance
func GetControllerInstance() *RelayController {
	return controllerInstance
}

func boolToInt(boolean bool) int {
	if boolean {
		return 1
	}

	return 0
}

func intToBool(integer int) bool {
	if integer == 0 {
		return false
	}

	return true
}
