package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
	"log"
	"os"
	"strconv"
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
	SwPin                 = 20
	VID            uint16 = 0x054C // SONY
	PID            uint16 = 0x06C1 // RC-S380
	Debug                 = true
)

func main() {
	log.Printf("%s /////// START OPEN KEY PROCESS ///////\n", DebugLogPrefix)

	initialize()
	OpenKey()
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

	//manageMosPin = rpio.Pin(MosPin) // MOS SEIGYO OUT PUT PIN
	//manageMosPin.Output()
	//manageMosPin.Low()
	managePWMPin = rpio.Pin(PwmPin) // SEIGYO OUT PUT PIN
	managePWMPin.Mode(rpio.Pwm)
	managePWMPin.Freq(50 * 100)
	managePWMPin.DutyCycle(0, 360)
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
	managePWMPin.High()

	go func() {
		time.Sleep(5000 * time.Millisecond)
		managePWMPin.Low()
		manageMosPin.Low()
	}()
	isOpenKey = true
	i := 0
	for {
		i++
		managePWMPin.DutyCycle(uint32(i), 360)
		fmt.Println("i " + strconv.Itoa(i))
		time.Sleep(100 * time.Millisecond)
		if manageSwPin.Read() == 1 {
			fmt.Println("-: -: end process " + strconv.Itoa(i))
			break
		}
	}

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
