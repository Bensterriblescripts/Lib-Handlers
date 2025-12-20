package osapi

import (
	"runtime"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

var LogKeys = true
var Hotkeys []Hotkey

type Hotkey struct {
	ID  uintptr
	Mod string
	Key string
}

var Keys = map[string]uintptr{
	"f1":  0x70,
	"f2":  0x71,
	"f3":  0x72,
	"f4":  0x73,
	"f5":  0x74,
	"f6":  0x75,
	"f7":  0x76,
	"f8":  0x77,
	"f9":  0x78,
	"f10": 0x79,
	"f11": 0x7A,
	"f12": 0x7B,
	"a":   0x41,
	"b":   0x42,
	"c":   0x43,
	"d":   0x44,
	"e":   0x45,
	"f":   0x46,
	"g":   0x47,
	"h":   0x48,
	"i":   0x49,
	"j":   0x4A,
	"k":   0x4B,
	"l":   0x4C,
	"m":   0x4D,
	"n":   0x4E,
	"o":   0x4F,
	"p":   0x50,
	"q":   0x51,
	"r":   0x52,
	"s":   0x53,
	"t":   0x54,
	"u":   0x55,
	"v":   0x56,
	"w":   0x57,
	"x":   0x58,
	"y":   0x59,
	"z":   0x5A,
	"0":   0x30,
	"1":   0x31,
	"2":   0x32,
	"3":   0x33,
	"4":   0x34,
	"5":   0x35,
	"6":   0x36,
	"7":   0x37,
	"8":   0x38,
	"9":   0x39,
}
var Modifiers = map[string]uintptr{
	"alt":     0x0001,
	"control": 0x0002,
	"shift":   0x0004,
	"win":     0x0008,
}

func StartKeylogger() {
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		Hotkeys = append(Hotkeys, Hotkey{ID: 1, Mod: "alt", Key: "f1"})
		for _, hotkey := range Hotkeys {
			if !RegisterHotKey(0, hotkey.ID, Modifiers[hotkey.Mod], Keys[hotkey.Key]) {
				ErrorLog("Failed to register hotkey: " + hotkey.Mod + "+" + hotkey.Key)
				return
			}
			defer UnregisterHotKey(0, hotkey.ID)
		}

		for LogKeys {
			var msg MSG
			r := GetMessage(&msg, 0, 0, 0)
			if r <= 0 {
				break
			}
			for _, hotkey := range Hotkeys {
				if msg.WParam == hotkey.ID {
					TraceLog("Hotkey Pressed: " + hotkey.Mod + " + " + hotkey.Key)
				}
			}
		}
	}()
}

func RegisterHotKey(hwnd uintptr, id uintptr, modifiers uintptr, vk uintptr) bool {
	r, _, _ := procRegisterHotKey.Call(
		hwnd,
		id,
		modifiers,
		vk,
	)
	return r != 0
}

func UnregisterHotKey(hwnd uintptr, id uintptr) bool {
	r, _, _ := procUnregisterHotKey.Call(
		hwnd,
		id,
	)
	return r != 0
}
