package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

var (
	managePWMPin    rpio.Pin
	manageSwPin     rpio.Pin
	manageSwDoorPin rpio.Pin
	position        int  = StopPosition
	motorRunning    bool = false
	motorStartTime  time.Time
)

const (
	DebugLogPrefix   = "[DEBUG]"
	PwmPin           = 13
	SwPin            = 21
	SwDoorPin        = 20
	StopPosition     = 1520 // サーボモーターを停止させるPWMパルス幅(マイクロ秒)
	ForwardPosition  = 800  // サーボモーターを正転させるPWMパルス幅(マイクロ秒)
	ReversePosition  = 2000 // サーボモーターを反転させるPWMパルス幅(マイクロ秒)
	IgnoreSwitchTime = 500  // スイッチ判定を無視する時間 (ミリ秒)
	timeout          = 1500 // 応答がなかった場合にタイムアウトして処理を終了する時間 (ミリ秒)
	Low              = 0
	High             = 1
)

func main() {
	log.Printf("%s /////// START OPEN KEY PROCESS ///////\n", DebugLogPrefix)

	fmt.Println("-: -: Servo setup...")
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	managePWMPin = rpio.Pin(PwmPin) // SEIGYO OUT PUT PIN
	managePWMPin.Mode(rpio.Pwm)
	managePWMPin.Freq(50 * 1000) // 50 Hz * DutyCycle (set servo)
	managePWMPin.DutyCycle(0, 1000)
	managePWMPin.High()

	// 歯車側
	manageSwPin = rpio.Pin(SwPin)
	manageSwPin.Input()
	manageSwPin.PullUp()

	// ドア側
	manageSwDoorPin = rpio.Pin(SwDoorPin)
	manageSwDoorPin.Input()
	manageSwDoorPin.PullUp()

	fmt.Println("-: -: end setup...")

	if len(os.Args) >= 2 {
		if os.Args[1] == "check" {
			for {
				if manageSwPin.Read() == rpio.Low {
					fmt.Println("sw pin push")
				}
				if manageSwDoorPin.Read() == rpio.Low {
					fmt.Println("sw door pin push")
				}
			}
		}
	}

	youyou()

	// // サーボモーターを制御
	// SetServo(managePWMPin, 2200)
	// time.Sleep(300 * time.Millisecond)
	// // SetServo(managePWMPin, 1700)
	// for {
	// 	if manageSwPin.Read() != 1 {
	// 		SetServo(managePWMPin, 0)
	// 		break
	// 	}
	// }

	// fmt.Println("Done 1")
	// SetServo(managePWMPin, 0)
	// time.Sleep(2500 * time.Millisecond)
	// SetServo(managePWMPin, 800)

	// // time.Sleep(300 * time.Millisecond)
	// SetServo(managePWMPin, 1300)
	// for {
	// 	if manageSwPin.Read() != 1 {
	// 		break
	// 	}
	// }
	// fmt.Println("Done 2")

	// SetServo(managePWMPin, 0)
}

// 指定したパルス幅でサーボモーターを制御
func SetServo(pin rpio.Pin, pulseWidthMicroSeconds float64) {
	// PWM周期（20ms）のうち、Highにすべき時間を計算
	pulseWidthFraction := pulseWidthMicroSeconds / 20000
	dutyCycle := uint32(pulseWidthFraction * 1000) // 分解能1000で計算

	pin.DutyCycle(dutyCycle, 1000)
}

func youyou() {
	fmt.Print("Enter command (f: forward, r: reverse, s: stop, e: exit): ")
	var command rune
	_, err := fmt.Scanf("%c\n", &command)
	if err != nil {
		fmt.Println("Error reading command:", err)
		youyou()
	}

	switch command {
	case 'f', 'r':
		if !motorRunning || time.Since(motorStartTime) > IgnoreSwitchTime*time.Millisecond {
			if command == 'f' {
				position = ForwardPosition
			} else {
				position = ReversePosition
			}
			motorStartTime = time.Now()
			motorRunning = true
		}
	case 's':
		position = StopPosition
		motorRunning = false
	case 'e':
		return
	}

	SetServo(managePWMPin, float64(position))

	for {

		if motorRunning && time.Since(motorStartTime) > IgnoreSwitchTime*time.Millisecond {
			if manageSwPin.Read() == rpio.Low {
				position = StopPosition
				motorRunning = false
				fmt.Printf("end servo\n")
				SetServo(managePWMPin, float64(position))
				break
			}
		}

		if motorRunning && time.Since(motorStartTime) > timeout*time.Millisecond {
			fmt.Printf("timeout...\n")
			position = StopPosition
			motorRunning = false
			SetServo(managePWMPin, float64(position))
			break
		}

		time.Sleep(20 * time.Millisecond) // ループの遅延

	}
	youyou()
}
