// *build windows

package vokiri

import (
	"strings"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
	"github.com/shu-go/rog"
)

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	enumWindows         = user32.NewProc("EnumWindows")
	enumChildWindows    = user32.NewProc("EnumChildWindows")
	getClassName        = user32.NewProc("GetClassNameW")
	getForegroundWindow = user32.NewProc("GetForegroundWindow")
	bringWindowToTop    = user32.NewProc("BringWindowToTop")
	getDlgItem          = user32.NewProc("GetDlgItem")
	getMenu             = user32.NewProc("GetMenu")
	getSubMenu          = user32.NewProc("GetSubMenu")
	getMenuItemID       = user32.NewProc("GetMenuItemID")
	getMenuItemCount    = user32.NewProc("GetMenuItemCount")
)

type EnumWindowsProc func(hwnd syscall.Handle, lparam uintptr) uintptr
type EnumChildProc func(hwnd syscall.Handle, lparam uintptr) uintptr

func EnumWindows(callback EnumWindowsProc, lparam uintptr) {
	enumWindows.Call(syscall.NewCallback(callback), lparam)
}

func EnumChildWindows(hwndParent syscall.Handle, callback EnumChildProc, lparam uintptr) {
	enumChildWindows.Call(uintptr(hwndParent), syscall.NewCallback(callback), lparam)
}

func GetDlgItem(hwnd syscall.Handle, id uintptr) syscall.Handle {
	ret, _, _ := getDlgItem.Call(uintptr(hwnd), id)
	return syscall.Handle(ret)
}

func GetWindowText(hwnd syscall.Handle) string {
	tlen := win.SendMessage(win.HWND(hwnd), win.WM_GETTEXTLENGTH, 0, 0)
	if tlen == 0 {
		return ""
	}
	tlen++

	buff := make([]uint16, tlen)
	win.SendMessage(win.HWND(hwnd), win.WM_GETTEXT, tlen, uintptr(unsafe.Pointer(&buff[0])))

	return syscall.UTF16ToString(buff)
}

func SetWindowText(hwnd syscall.Handle, s string) {
	win.SendMessage(win.HWND(hwnd), win.WM_SETTEXT, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s))))
}

func SetFocus(hwnd syscall.Handle) {
	win.SendMessage(win.HWND(hwnd), win.WM_SETFOCUS, 0, 0)
}

func LoseFocus(hwnd syscall.Handle) {
	win.SendMessage(win.HWND(hwnd), win.WM_KILLFOCUS, 0, 0)
}

func SendChar(hwnd syscall.Handle, c byte) {
	win.SendMessage(win.HWND(hwnd), win.WM_CHAR, uintptr(c), 0)
}

func GetClassName(hwnd syscall.Handle) string {
	const buffLen = 100
	buff := make([]uint16, buffLen)

	getClassName.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buff[0])), buffLen)

	return syscall.UTF16ToString(buff)
}

func GetForegroundWindow() syscall.Handle {
	w, _, _ := getForegroundWindow.Call()
	return syscall.Handle(w)
}

func BringWindowToTop(hwnd syscall.Handle) {
	bringWindowToTop.Call(uintptr(hwnd))
}

func GetMenu(hwnd syscall.Handle) syscall.Handle {
	m, _, _ := getMenu.Call(uintptr(hwnd))
	return syscall.Handle(m)
}

func GetSubMenu(hmenu syscall.Handle, pos int32) syscall.Handle {
	m, _, _ := getSubMenu.Call(uintptr(hmenu), uintptr(pos))
	return syscall.Handle(m)
}

func GetMenuItemID(hmenu syscall.Handle, pos int32) int32 {
	id, _, _ := getMenuItemID.Call(uintptr(hmenu), uintptr(pos))
	return int32(id)
}

func GetMenuItemCount(hmenu syscall.Handle) int32 {
	c, _, _ := getMenuItemCount.Call(uintptr(hmenu))
	return int32(c)
}

func SetRichEditText(hwnd syscall.Handle, s string) {
	win.SendMessage(win.HWND(hwnd), win.EM_SETSEL, 0, 0xffffffff)
	win.SendMessage(win.HWND(hwnd), win.EM_REPLACESEL, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s))))
}

func PushButton(hwnd syscall.Handle) {
	win.PostMessage(win.HWND(hwnd), win.BM_CLICK, 0, 0)
}

func ChangeTab(hwnd syscall.Handle, index uintptr) {
	i := win.SendMessage(win.HWND(hwnd), win.TCM_SETCURSEL, index, 0)
	if i == 0xffffffff {
		rog.Debug("TCM_SETCURSEL", i)
	}
	i = win.SendMessage(win.HWND(hwnd), win.TCM_SETCURFOCUS, index, 0)
	if i == 0xffffffff {
		rog.Debug("TCM_SETCURFOCUS", i)
	}
}

func FindWindow(patterns string) syscall.Handle {
	var handle syscall.Handle

	EnumWindows(func(hwnd syscall.Handle, _ uintptr) uintptr {
		title := strings.ToUpper(GetWindowText(hwnd))

		matches := true
		for _, p := range strings.Split(patterns, " ") {
			if !strings.Contains(title, strings.ToUpper(p)) {
				matches = false
				break
			}
		}
		if matches {
			//rog.Debug(title)
			handle = hwnd
			return 0
		}

		return 1
	}, 0)

	return handle
}

func FindChildWindow(hparent syscall.Handle, c, t string, exact bool) syscall.Handle {
	var handle syscall.Handle

	c = strings.ToUpper(c)
	t = strings.ToUpper(t)

	if len(t) == 0 && len(c) == 0 {
		return handle
	}

	EnumChildWindows(hparent, func(hwnd syscall.Handle, _ uintptr) uintptr {
		text := strings.ToUpper(GetWindowText(hwnd))
		class := strings.ToUpper(GetClassName(hwnd))

		if len(c) > 0 && (exact && class == c || !exact && strings.Contains(class, c)) &&
			len(t) > 0 && (exact && text == t || !exact && strings.Contains(text, t)) {
			//
			handle = hwnd
			return 0
		}

		return 1
	}, 0)

	return handle
}

func FindChildWindowByClass(hparent syscall.Handle, c string, exact bool) syscall.Handle {
	var handle syscall.Handle

	c = strings.ToUpper(c)

	EnumChildWindows(hparent, func(hwnd syscall.Handle, _ uintptr) uintptr {
		class := strings.ToUpper(GetClassName(hwnd))

		if exact && class == c || !exact && strings.Contains(class, c) {
			handle = hwnd
			return 0
		}

		return 1
	}, 0)

	return handle
}

func FindChildWindowByText(hparent syscall.Handle, t string, exact bool) syscall.Handle {
	var handle syscall.Handle

	t = strings.ToUpper(t)

	EnumChildWindows(hparent, func(hwnd syscall.Handle, _ uintptr) uintptr {
		text := strings.ToUpper(GetWindowText(hwnd))

		if exact && text == t || !exact && strings.Contains(text, t) {
			handle = hwnd
			return 0
		}

		return 1
	}, 0)

	return handle
}

func Sendkeys(keys []win.KEYBD_INPUT) uint32 {
	return win.SendInput(uint32(len(keys)), unsafe.Pointer(&keys), int32(unsafe.Sizeof(keys[0])))
}

func Sendkey(key win.KEYBDINPUT) uint32 {
	return Sendkeys([]win.KEYBD_INPUT{{Type: win.INPUT_KEYBOARD, Ki: key}})
}

func ShowWindow(hwnd syscall.Handle, show bool) {
	if show {
		win.ShowWindow(win.HWND(hwnd), win.SW_SHOW)
	} else {
		win.ShowWindow(win.HWND(hwnd), win.SW_HIDE)
	}
}
