package main

import (
	"fmt"
	"github.com/bamchoh/pasori"
	"log"
	"os"
	"time"
)

var (
	managePWMPin    rpio.Pin
	manageMosPin    rpio.Pin
	manageSwPin     rpio.Pin
	isOpenKey       bool
	lock            bool   = false
	isRegister      bool   = false
	isCloseProgress bool   = false
	tempName        string = ""
)

const (
	DebugLogPrefix        = "[DEBUG]"
	PwmPin                = 13
	MosPin                = 17
	SwPin                 = 18
	VID            uint16 = 0x054C // SONY
	PID            uint16 = 0x06C1 // RC-S380
	Debug                 = true
)

func main() {
	log.Printf("%s /////// START OPEN KEY PROCESS ///////\n", DebugLogPrefix)

	initialize()

	for {
		// sudoしないと動かないので注意
		idm, err := pasori.GetID(VID, PID)
		if err != nil {
			continue
		}

		log.Println(idm)

		if isOpenKey {
			CloseKey()
		} else {
			OpenKey()
		}

		time.Sleep(2000 * time.Millisecond)
	}
}

func initialize() {
	log.Printf("%s -: Initializing -----\n", DebugLogPrefix)

	////////////////// SERVO

	_ = os.MkdirAll("data", 0755)

	fmt.Println("-: -: Servo setup...")
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	manageMosPin = rpio.Pin(MosPin) // MOS SEIGYO OUT PUT PIN
	manageMosPin.Output()
	manageMosPin.Low()
	managePWMPin = rpio.Pin(PwmPin) // SEIGYO OUT PUT PIN
	managePWMPin.Mode(rpio.Pwm)
	managePWMPin.Freq(50 * 100)
	managePWMPin.DutyCycle(0, 100)
	managePWMPin.Low()
	fmt.Println("-: -: END Servo setup")

	////////////////// SWITCH
	fmt.Println("-: -: switch setup...")

	manageSwPin = rpio.Pin(SwPin)
	manageSwPin.Input()
	manageSwPin.PullUp()
	//manageSwPin.Detect(rpio.FallEdge)

	////////////////// PASORI
	fmt.Println("-: -: IDM Read setup...")

}

func OpenKey() {
	lock = true
	manageMosPin.High()
	managePWMPin.High()

	time.Sleep(500 * time.Millisecond)

	for i := 1; i <= 60; i++ {
		managePWMPin.DutyCycle(uint32(i), 100)
		time.Sleep(10 * time.Millisecond)
	}

	go func() {
		time.Sleep(1000 * time.Millisecond)
		managePWMPin.Low()
		manageMosPin.Low()
	}()
	isOpenKey = true

}

func CloseKey() {
	manageMosPin.High()
	managePWMPin.High()
	time.Sleep(500 * time.Millisecond)
	for i := 1; i <= 60; i++ {
		managePWMPin.DutyCycle(uint32(50-i), 100)
		time.Sleep(10 * time.Millisecond)
	}
	go func() {
		time.Sleep(1000 * time.Millisecond)
		managePWMPin.Low()
		manageMosPin.Low()
	}()
	isOpenKey = false
}
