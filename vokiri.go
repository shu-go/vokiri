// *build windows

package vokiri

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	//"unicode"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"bitbucket.org/shu/log"
	"bitbucket.org/shu/retry"
	"github.com/lxn/win"
)

type Kiritan struct {
	ExePath     string
	WindowTitle string

	PDictPath string // if injected

	WindowHandle syscall.Handle
	REditHandle  syscall.Handle

	PlayButtonHandle syscall.Handle
	StopButtonHandle syscall.Handle
	SaveButtonHandle syscall.Handle

	TabHandle          syscall.Handle
	VolumeEditHandle   syscall.Handle
	SpeedEditHandle    syscall.Handle
	PitchEditHandle    syscall.Handle
	EmphasisEditHandle syscall.Handle
}

func New(path, title string) *Kiritan {
	v := &Kiritan{
		ExePath:     path,
		WindowTitle: title,
	}
	v.Rescan()
	return v
}

func (v *Kiritan) Rescan() {
	v.WindowHandle = 0
	v.REditHandle = 0
	v.PlayButtonHandle = 0
	v.StopButtonHandle = 0
	v.SaveButtonHandle = 0
	v.VolumeEditHandle = 0
	v.SpeedEditHandle = 0
	v.PitchEditHandle = 0
	v.EmphasisEditHandle = 0
	v.TabHandle = 0

	log.Infof("VOICEROIDのウィンドウを確認...(%q)", v.WindowTitle)
	v.WindowHandle = FindWindow(v.WindowTitle)
	if v.WindowHandle == 0 {
		log.Info("  現時点では見つかりませんでした。今回のスキャンはこれで打ち切ります。")
		return
	} else {
		log.Infof("  ウィンドウハンドル=%X", v.WindowHandle)
		log.Info("  起動が確認できました。")
	}

	log.Info("ウィンドウ中のコントロールから、操作に必要なものを検索...")

	EnumChildWindows(v.WindowHandle, func(hwnd syscall.Handle, _ uintptr) uintptr {
		title := GetWindowText(hwnd)
		class := strings.ToUpper(GetClassName(hwnd))

		//log.Debug(hwnd, "Class:", class, "Title:", title)

		if v.REditHandle == 0 && strings.Contains(class, "RICHEDIT") {
			// first RichEdit
			v.REditHandle = hwnd
		} else if v.PlayButtonHandle == 0 && strings.Contains(class, "BUTTON") && (strings.Contains(title, "再生") && !strings.Contains(title, "時間") || strings.Contains(title, "一時停止")) {
			v.PlayButtonHandle = hwnd
		} else if v.StopButtonHandle == 0 && strings.Contains(class, "BUTTON") && strings.Contains(title, "停止") {
			v.StopButtonHandle = hwnd
		} else if v.SaveButtonHandle == 0 && strings.Contains(class, "BUTTON") && strings.Contains(title, "音声保存") {
			v.SaveButtonHandle = hwnd
		} else if /* last -- v.TabHandle == 0 &&*/ strings.Contains(class, "SYSTABCONTROL") {
			v.TabHandle = hwnd
		}

		return 1
	}, 0)

	if v.PlayButtonHandle != 0 {
		log.Infof("  再生ボタン=%X", v.PlayButtonHandle)
		log.Info("  再生ボタンが見つかりました。")
	}
	if v.SaveButtonHandle != 0 {
		log.Infof("  音声保存ボタン=%x", v.SaveButtonHandle)
		log.Info("  音声保存ボタンが見つかりました。")
	}

	log.Infof("(オプション)音声効果をふくむタブ=%X", v.TabHandle)

	if v.TabHandle != 0 {
		log.Info("音声効果をふくむタブが見つかったので、各種パラメーターのテキストボックスを検索...")

		var edits [4]syscall.Handle
		ei := 0

		//ChangeTab(v.TabHandle, 2)
		x := 200
		y := 10
		win.SendMessage(win.HWND(v.TabHandle), win.WM_LBUTTONDOWN, win.MK_LBUTTON, uintptr(y<<16|x))
		win.SendMessage(win.HWND(v.TabHandle), win.WM_LBUTTONUP, 0, 0)

		EnumChildWindows(v.WindowHandle, func(hwnd syscall.Handle, _ uintptr) uintptr {
			class := GetClassName(hwnd)

			if strings.Contains(class, "EDIT") && ei < len(edits) {
				edits[ei] = hwnd
				ei++
				if ei == len(edits) {
					return 0
				}
			} else {
				ei = 0
			}

			return 1
		}, 0)

		v.VolumeEditHandle = edits[3]
		v.SpeedEditHandle = edits[2]
		v.PitchEditHandle = edits[1]
		v.EmphasisEditHandle = edits[0]

		log.Infof("  音量=%X", v.VolumeEditHandle)
		log.Infof("  話速=%X", v.SpeedEditHandle)
		log.Infof("  高さ=%X", v.PitchEditHandle)
		log.Infof("  抑揚=%X", v.EmphasisEditHandle)

		if v.VolumeEditHandle != 0 && v.SpeedEditHandle != 0 && v.PitchEditHandle != 0 && v.EmphasisEditHandle != 0 {
			log.Info("  音声効果の各種パラメーターの指定が利用可能です。")
		} else {
			log.Info("  音声効果の各種パラメーターのテキストボックスが見つかりませんでした。＊現時点では＊利用できません。")
		}
	}
}

func (v *Kiritan) Close() error {
	if !v.IsRunning() {
		return nil
	}

	log.Info("終了操作開始...")

	log.Info("  操作可能になるまで待機")
	v.WaitForSpeech()

	log.Info("  終了メッセージ送信")
	win.PostMessage(win.HWND(v.WindowHandle), win.WM_CLOSE, 0, 0)

	log.Info("  終了の確認ダイアログボックスが出ていないか確認...")
	var hmsg syscall.Handle
	retry.Wait(3*time.Second, WAIT_SHORT, func() bool {
		hmsg = FindWindow("注意")
		if hmsg != 0 {
			return true
		}
		return false
	})
	if hmsg != 0 {
		yesButton := FindChildWindowByClass(hmsg, "BUTTON", false)
		log.Infof("      はいボタン=%X", yesButton)
		if yesButton != 0 {
			win.SendMessage(win.HWND(yesButton), win.BM_CLICK, 0, 0)
		}
	}
	return nil
}

func (v *Kiritan) IsRunning() bool {
	return v.WindowHandle != 0 && v.PlayButtonHandle != 0 && v.SaveButtonHandle != 0
}

func (v *Kiritan) Run() error {
	cmd := exec.Command(v.ExePath)
	return cmd.Start()
}

func (v *Kiritan) WaitForRun() bool {
	log.Info("VOICEROIDの起動を待っています。")

	done := retry.Wait(20*time.Second, WAIT_LONG, func() bool {
		if v.IsRunning() {
			return true
		}
		v.Rescan()
		return false
	})

	if done {
		return true
	}

	log.Info("規定時間内に起動を確認できませんでした。")

	return false
}

func (v *Kiritan) waitForForeground() bool {
	if v.WindowHandle == 0 {
		log.Info("VOICEROIDが起動していません。")
		return false
	}

	BringWindowToTop(v.WindowHandle)

	first := true
	done := retry.Wait(time.Second, WAIT_SHORT, func() bool {
		hwnd := GetForegroundWindow()
		if hwnd == v.WindowHandle {
			return true
		}

		if first {
			first = false
			v.Run()
		}
		return false
	})
	win.SetFocus(win.HWND(v.WindowHandle))

	if done {
		return true
	}

	return false
}

func (v *Kiritan) WaitForSpeech() bool {
	return v.waitForEnabled(v.SaveButtonHandle, true, 60*time.Second, WAIT_LONG)
}

func (v *Kiritan) waitForStartSpeech() bool {
	return v.waitForEnabled(v.SaveButtonHandle, false, 60*time.Second, WAIT_LONG)
}

func (v *Kiritan) waitForEnabled(hwnd syscall.Handle, wants bool, timeout, wait time.Duration) bool {
	if v.WindowHandle == 0 {
		return false
	}

	done := retry.Wait(timeout, wait, func() bool {
		enabled := win.IsWindowEnabled(win.HWND(hwnd))
		if wants && enabled || !wants && !enabled {
			return true
		}
		return false
	})
	return done
}

func (v *Kiritan) saveToFile(dest string) error {
	log.Info("保存操作開始...")

	var hdlg syscall.Handle
	retry.Wait(5*time.Second, WAIT_SHORT, func() bool {
		hdlg = FindWindow("音声ファイルの保存")
		if hdlg != 0 {
			return true
		}
		return false
	})
	if hdlg == 0 {
		return fmt.Errorf("音声ファイルの保存ダイアログボックスが見つかりませんでした。")
	}

	var filenameEdit syscall.Handle
	retry.Wait(5*time.Second, WAIT_SHORT, func() bool {
		filenameEdit = FindChildWindowByClass(hdlg, "Edit", true)
		if filenameEdit != 0 {
			return true
		}
		return false
	})
	if filenameEdit == 0 {
		return fmt.Errorf("音声ファイルの保存ダイアログボックスで、ファイル名テキストボックスが見つかりませんでした。")
	}

	log.Infof("  ファイル名テキストボックス=%X", filenameEdit)
	log.Infof("    内容=%s", dest)

	SetWindowText(filenameEdit, dest)
	retry.Wait(1*time.Second, WAIT_SHORT, func() bool {
		d := GetWindowText(filenameEdit)
		if d == dest {
			return true
		}
		SetWindowText(filenameEdit, dest)
		return false
	})

	defButtonID := win.SendMessage(win.HWND(hdlg) /*DM_GETDEFID*/, (win.WM_USER+0), 0, 0) & 0xffff
	if defButtonID == 0 {
		return fmt.Errorf("音声ファイルの保存ダイアログボックスで、保存ボタンが見つかりませんでした。")
	}
	defButton := GetDlgItem(hdlg, defButtonID)
	log.Infof("  保存ボタン=%X", defButton)
	if defButton == 0 {
		return fmt.Errorf("音声ファイルの保存ダイアログボックスで、保存ボタンが見つかりませんでした。")
	}

	log.Info("  保存ボタン押下")
	win.PostMessage(win.HWND(defButton), win.BM_CLICK, 0, 0)

	log.Info("  上書き保存の確認ダイアログボックスが出ていないか確認...")
	var hmsg syscall.Handle
	first := true
	retry.Wait(time.Second, WAIT_SHORT, func() bool {
		hmsg = FindWindow("保存 確認")
		if hmsg != 0 {
			return true
		}
		if !first {
			hmsg = FindWindow("音声ファイルの保存")
			if hmsg != 0 {
				return true
			}
		}
		first = false
		return false
	})
	log.Infof("    確認ダイアログボックス=%X", hmsg)
	if hmsg != 0 {
		yesButton := FindChildWindowByClass(hmsg, "BUTTON", false)
		log.Infof("      はいボタン=%X", yesButton)
		if yesButton != 0 {
			win.SendMessage(win.HWND(yesButton), win.BM_CLICK, 0, 0)
		}
	}

	retry.Wait(5*time.Second, WAIT_VERYSHORT, func() bool {
		hdlg = FindWindow("音声ファイルの保存")
		if hdlg == 0 {
			return true
		}
		return false
	})

	return nil
}

func (v *Kiritan) SetVolume(vol float64) error {
	log.Infof("音量->%.2f", vol)

	if v.VolumeEditHandle == 0 {
		return fmt.Errorf("音量テキストボックスが見つかりませんでした。")
	}

	return v.setParameterEditBox(v.VolumeEditHandle, vol)
}

func (v *Kiritan) SetSpeed(s float64) error {
	log.Infof("話速->%.2f", s)

	if v.SpeedEditHandle == 0 {
		return fmt.Errorf("話速テキストボックスが見つかりませんでした。")
	}

	return v.setParameterEditBox(v.SpeedEditHandle, s)
}

func (v *Kiritan) SetPitch(p float64) error {
	log.Infof("高さ->%.2f", p)

	if v.PitchEditHandle == 0 {
		return fmt.Errorf("高さテキストボックスが見つかりませんでした。")
	}

	return v.setParameterEditBox(v.PitchEditHandle, p)
}

func (v *Kiritan) SetEmphasis(i float64) error {
	log.Infof("抑揚->%.2f", i)

	if v.EmphasisEditHandle == 0 {
		return fmt.Errorf("抑揚テキストボックスが見つかりませんでした。")
	}

	return v.setParameterEditBox(v.EmphasisEditHandle, i)
}

func (v *Kiritan) setParameterEditBox(hwnd syscall.Handle, value float64) error {
	//v.waitForForeground()

	v.WaitForSpeech()

	v.waitForEnabled(hwnd, true, 5*time.Second, WAIT_SHORT)

	SetWindowText(hwnd, fmt.Sprintf("%.2f", value))

	win.SendMessage(win.HWND(hwnd), win.WM_LBUTTONDOWN, win.MK_LBUTTON, 0)
	win.SendMessage(win.HWND(hwnd), win.WM_LBUTTONUP, 0, 0)
	win.SendMessage(win.HWND(v.REditHandle), win.WM_LBUTTONDOWN, win.MK_LBUTTON, 0)
	win.SendMessage(win.HWND(v.REditHandle), win.WM_LBUTTONUP, 0, 0)

	return nil
}

func (v *Kiritan) GetVolume() float64 {
	return v.getParameterEditBox(v.VolumeEditHandle)
}

func (v *Kiritan) GetSpeed() float64 {
	return v.getParameterEditBox(v.SpeedEditHandle)
}

func (v *Kiritan) GetPitch() float64 {
	return v.getParameterEditBox(v.PitchEditHandle)
}

func (v *Kiritan) GetEmphasis() float64 {
	return v.getParameterEditBox(v.EmphasisEditHandle)
}

func (v *Kiritan) getParameterEditBox(hwnd syscall.Handle) float64 {
	s := GetWindowText(hwnd)
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 1.0
	}
	return value
}

func (v *Kiritan) Speak(s string) error {
	log.Infof("読み上げ「%s」", s)

	if v.REditHandle == 0 {
		return fmt.Errorf("テキストボックスが見つかりませんでした。")
	}
	if v.PlayButtonHandle == 0 {
		return fmt.Errorf("再生ボタンが見つかりませんでした。")
	}
	if v.StopButtonHandle == 0 {
		return fmt.Errorf("停止ボタンが見つかりませんでした。")
	}
	if v.SaveButtonHandle == 0 {
		return fmt.Errorf("音声保存ボタンが見つかりませんでした。")
	}

	v.WaitForSpeech()

	SetRichEditText(v.REditHandle, s)
	PushButton(v.PlayButtonHandle)

	time.Sleep(WAIT_SHORT)

	//log.Debug(v.PDictPath)
	if len(v.PDictPath) > 0 {
		v.InjectPhraseEntries(nil)
	}

	return nil
}

func (v *Kiritan) Record(s, dest string) error {
	log.Infof("録音「%s」=> %q", s, dest)

	if v.REditHandle == 0 {
		return fmt.Errorf("テキストボックスが見つかりませんでした。")
	}
	if v.PlayButtonHandle == 0 {
		return fmt.Errorf("再生ボタンが見つかりませんでした。")
	}
	if v.StopButtonHandle == 0 {
		return fmt.Errorf("停止ボタンが見つかりませんでした。")
	}
	if v.SaveButtonHandle == 0 {
		return fmt.Errorf("音声保存ボタンが見つかりませんでした。")
	}

	v.WaitForSpeech()

	SetRichEditText(v.REditHandle, s)
	PushButton(v.SaveButtonHandle)

	v.saveToFile(dest)

	if len(v.PDictPath) > 0 {
		v.InjectPhraseEntries(nil)
	}

	return nil
}

func (v *Kiritan) ReloadPhraseDict(pEntries []string, persist bool) error {
	v.waitForForeground()
	v.WaitForSpeech()

	win.PostMessage(win.HWND(v.WindowHandle), win.WM_SYSKEYDOWN, 'F', 1<<29|1)
	time.Sleep(time.Millisecond)
	win.PostMessage(win.HWND(v.WindowHandle), win.WM_SYSKEYDOWN, win.VK_RIGHT, 1<<29|1)
	time.Sleep(time.Millisecond)
	win.PostMessage(win.HWND(v.WindowHandle), win.WM_SYSKEYDOWN, 'S', 1<<29|1)
	time.Sleep(time.Millisecond)
	win.PostMessage(win.HWND(v.WindowHandle), win.WM_SYSKEYDOWN, 'L', 1<<29|1)

	var hdlg syscall.Handle
	retry.Wait(100*time.Millisecond, WAIT_VERYSHORT, func() bool {
		hdlg = FindWindow("日本語辞書設定")
		return hdlg != 0
	})
	if hdlg == 0 {
		return fmt.Errorf("日本語辞書設定ダイアログボックスが開けませんでした。")
	}

	var pdictEdit syscall.Handle
	EnumChildWindows(hdlg, func(hwnd syscall.Handle, _ uintptr) uintptr {
		title := strings.ToUpper(GetWindowText(hwnd))
		class := strings.ToUpper(GetClassName(hwnd))

		//log.Debug(hwnd, class, title)
		if strings.Contains(title, ".PDIC") && strings.Contains(class, "EDIT") {
			pdictEdit = hwnd
			return 0
		}

		return 1
	}, 0)
	if pdictEdit == 0 {
		return fmt.Errorf("日本語辞書設定ダイアログボックスで、辞書ファイルの指定テキストボックスが見つかりませんでした。。")
	}

	v.PDictPath = GetWindowText(pdictEdit)
	if persist {
		if err := v.InjectPersistedPhraseEntries(pEntries); err != nil {
			return fmt.Errorf("フレーズ辞書の編集に失敗しました: %v", err)
		}
	} else {
		if err := v.InjectPhraseEntries(pEntries); err != nil {
			return fmt.Errorf("フレーズ辞書の編集に失敗しました: %v", err)
		}
	}

	cb := FindChildWindowByText(hdlg, "フレーズ辞書", false)
	if cb == 0 {
		log.Info("フレーズ辞書のチェックボックスが見つかりませんでした。")
	} else {
		win.SendMessage(win.HWND(cb), win.BM_SETCHECK, win.BST_CHECKED, 0)
		retry.Wait(100*time.Millisecond, WAIT_VERYSHORT, func() bool {
			return win.SendMessage(win.HWND(cb), win.BM_GETCHECK, 0, 0) == win.BST_CHECKED
		})
	}

	okButton := FindChildWindow(hdlg, "BUTTON", "OK", false)
	if okButton == 0 {
		return fmt.Errorf("日本語辞書設定ダイアログボックスで、OKボタンが見つかりませんでした。。")
	}
	win.SendMessage(win.HWND(okButton), win.BM_CLICK, 0, 0)

	// wait for closed
	retry.Wait(time.Second, WAIT_SHORT, func() bool {
		hdlg = FindWindow("日本語辞書設定")
		return hdlg == 0
	})

	return nil
}

func (v *Kiritan) InjectPhraseEntries(pEntries []string) error {
	lines := make([]string, 0, 20)

	f, err := os.Open(v.PDictPath)
	if err != nil {
		f.Close()
		return fmt.Errorf("フレーズ辞書の読み込みに失敗しました: %v", err)
	}

	scanner := bufio.NewScanner(transform.NewReader(bufio.NewReader(f), japanese.ShiftJIS.NewDecoder()))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		f.Close()
		return fmt.Errorf("フレーズ辞書の読み込みに失敗しました: %v", err)
	}
	f.Close()

	// PDictPath will be written later

	// drop entries num>=1000
	pos := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "num:") {
			if num, err := strconv.Atoi(line[4:]); err != nil {
				return fmt.Errorf("フレーズ辞書の番号解析に失敗しました: %v", err)
			} else if num >= 1000 {
				pos = i
				break
			}
		}
	}
	if pos != -1 {
		lines = lines[:pos]
	}

	buff := &bytes.Buffer{}
	for _, line := range lines {
		buff.Write([]byte(line))
		buff.Write([]byte("\r\n"))
	}

	// append entries
	for i, e := range pEntries {
		pair := strings.Split(e, "=")
		if len(pair) < 2 {
			return fmt.Errorf("フレーズの区切りに失敗しました: %q", e)
		}
		s, p := pair[0], pair[1]
		if p != "" {
			//log.Debug(fmt.Sprintf("num:%d\r\n%s\r\n$2_2%s$2_2\r\n", 1000+i, s, v.CompilePhonation(p)))
			buff.Write([]byte(fmt.Sprintf("num:%d\r\n%s\r\n$2_2%s$2_2\r\n", 1000+i, s, v.CompilePhonation(p))))
		}
	}

	if data, err := ioutil.ReadAll(transform.NewReader(bufio.NewReader(buff), japanese.ShiftJIS.NewEncoder())); err != nil {
		return fmt.Errorf("内部変換処理に失敗しました: %v", err)
	} else if err := ioutil.WriteFile(v.PDictPath, data, os.ModePerm); err != nil {
		return fmt.Errorf("フレーズ辞書の保存に失敗しました: %v", err)
	}

	return nil
}

var (
	instPhraseStart = []string{"((", "（（"}
	instPhraseSep   = []string{"|", "｜"}
	instPhraseEnd   = []string{"))", "））"}
)

func ExtractInstantPhrases(text string) (pEntries []string, after string, err error) {
	var start int
	for {
		//log.Debug("text", text)
		//log.Debug("start", start)
		spos, slen := indexAnyAry(text[start:], instPhraseStart)
		pos := start + spos
		if pos < start {
			break
		}
		//log.Debug("pos", pos)

		epos, elen := indexAnyAry(text[pos:], instPhraseEnd)
		endpos := pos + epos
		if endpos < pos {
			break
		}
		//log.Debug("endpos", endpos)

		seppos, seplen := indexAnyAry(text[pos:endpos], instPhraseSep)
		seppos = pos + seppos
		if seppos < pos {
			start = endpos
			continue
		}
		//log.Debug("seppos", seppos)

		phrase := text[pos+slen : seppos]
		pronun := text[seppos+seplen : endpos]
		//log.Debug("phrase, pronun", phrase, pronun)
		pEntries = append(pEntries, strings.Trim(phrase, " ")+"="+strings.Trim(pronun, " "))

		start = pos + len(phrase) + 2
		text = text[:pos] + " " + phrase + " " + text[endpos+elen:]

		if len(text) <= start {
			break
		}
	}

	return pEntries, text, nil
}

func (v *Kiritan) InjectPersistedPhraseEntries(pEntries []string) error {
	prelines := make([]string, 0, 10)
	postlines := make([]string, 0, 10)
	type txt2pDictEntry struct {
		Num       string
		Text      string
		Phonation string
	}
	txt2pDict := make(map[string]txt2pDictEntry)

	f, err := os.Open(v.PDictPath)
	if err != nil {
		f.Close()
		return fmt.Errorf("フレーズ辞書の読み込みに失敗しました: %v", err)
	}

	part := 1
	scanner := bufio.NewScanner(transform.NewReader(bufio.NewReader(f), japanese.ShiftJIS.NewDecoder()))
	for scanner.Scan() {
		line := scanner.Text()

		if part == 1 {
			if strings.HasPrefix(line, "num:") {
				part = 2
			} else {
				prelines = append(prelines, scanner.Text())
			}
		}

		if part == 2 {
			if strings.HasPrefix(line, "num:") {
				if !scanner.Scan() {
					break
				}
				line2 := scanner.Text()

				if !scanner.Scan() {
					break
				}
				line3 := scanner.Text()

				txt2pDict[line2] = txt2pDictEntry{
					Num:       line,
					Text:      line2,
					Phonation: line3,
				}
			} else {
				part = 3
			}
		}

		if part == 3 {
			postlines = append(postlines, scanner.Text())
		}

	}
	if err := scanner.Err(); err != nil {
		f.Close()
		return fmt.Errorf("フレーズ辞書の読み込みに失敗しました: %v", err)
	}
	f.Close()

	// PDictPath will be written later

	// append entries
	for _, e := range pEntries {
		pair := strings.Split(e, "=")
		if len(pair) < 2 {
			return fmt.Errorf("フレーズの区切りに失敗しました: %q", e)
		}
		s, p := pair[0], pair[1]
		if p == "" {
			delete(txt2pDict, s)
		} else {
			if old, found := txt2pDict[s]; found {
				old.Phonation = fmt.Sprintf("$2_2%s$2_2", v.CompilePhonation(p))
				txt2pDict[s] = old
			} else {
				txt2pDict[s] = txt2pDictEntry{
					Num:       fmt.Sprintf("num:%d", len(txt2pDict)),
					Text:      s,
					Phonation: fmt.Sprintf("$2_2%s$2_2", v.CompilePhonation(p)),
				}
			}
		}
	}

	buff := &bytes.Buffer{}
	for _, line := range prelines {
		buff.Write([]byte(line))
		buff.Write([]byte("\r\n"))
	}

	sorted := make([]txt2pDictEntry, 0, len(txt2pDict))
	for _, v := range txt2pDict {
		sorted = append(sorted, v)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i].Num) < len(sorted[j].Num) || sorted[i].Num < sorted[j].Num
	})
	for _, v := range sorted {
		buff.Write([]byte(v.Num))
		buff.Write([]byte("\r\n"))
		buff.Write([]byte(v.Text))
		buff.Write([]byte("\r\n"))
		buff.Write([]byte(v.Phonation))
		buff.Write([]byte("\r\n"))
	}

	for _, line := range postlines {
		buff.Write([]byte(line))
		buff.Write([]byte("\r\n"))
	}

	if data, err := ioutil.ReadAll(transform.NewReader(bufio.NewReader(buff), japanese.ShiftJIS.NewEncoder())); err != nil {
		return fmt.Errorf("内部変換処理に失敗しました: %v", err)
	} else if err := ioutil.WriteFile(v.PDictPath, data, os.ModePerm); err != nil {
		return fmt.Errorf("フレーズ辞書の保存に失敗しました: %v", err)
	}

	return nil
}

var vol *regexp.Regexp
var spd *regexp.Regexp
var pit *regexp.Regexp
var emph *regexp.Regexp

func init() {
	vol = regexp.MustCompile(`\((?i:v|vol):?\s?([0-2](\.\d+)?)\)`)
	spd = regexp.MustCompile(`\((?i:s|spd):?\s?([0-2](\.\d+)?)\)`)
	pit = regexp.MustCompile(`\((?i:p|pit):?\s?([0-2](\.\d+)?)\)`)
	emph = regexp.MustCompile(`\((?i:e|emph):?\s?([0-2](\.\d+)?)\)`)
}

func (v *Kiritan) CompilePhonation(p string) string {
	//log.Debug(p)
	p = vol.ReplaceAllString(p, "(Vol ABSLEVEL=$1)")
	p = spd.ReplaceAllString(p, "(Spd ABSSPEED=$1)")
	p = pit.ReplaceAllString(p, "(Pit ABSLEVEL=$1)")
	p = emph.ReplaceAllString(p, "(EMPH ABSLEVEL=$1)")
	//log.Debug(p)
	return p
}

func (v *Kiritan) Test() {
	log.Debug("==================================================")

	text := "皆さん、（（こんにちは｜コ!ンニチ|ワー））"
	if pEntries, after, err := ExtractInstantPhrases(text); err != nil {
		log.Debug(err)
		return
	} else {
		log.Debug(text)
		log.Debug("  ", after)
		log.Debug("  ", pEntries)
	}
}
