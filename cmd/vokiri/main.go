// *build windows

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/andrew-d/go-termutil"
	"github.com/urfave/cli"

	"bitbucket.org/shu_go/rog"
	"bitbucket.org/shu_go/vokiri"
)

func main() {
	app := cli.NewApp()
	app.Name = "vokiri"
	app.Usage = "VOICEROID＋ 東北きりたん をコマンドラインから操作します"
	app.Version = "0.2.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "title", Value: "VOICEROID＋ 東北きりたん EX", Usage: "VOICEROIDのウィンドウタイトル(操作対象とみなすための条件)"},
		cli.StringFlag{Name: "exe", Value: "", Usage: "VOICEROIDの実行パス"},
		cli.BoolFlag{Name: "close", Usage: "実行後にVOICEROIDのウィンドウを終了させます"},
		cli.StringFlag{Name: "record, r", Value: "", Usage: "保存先のWAVファイル名"},
		cli.StringFlag{Name: "record-once", Value: "", Usage: "保存先のWAVファイル名(同名が存在する場合は何もしない)"},
		cli.Float64Flag{Name: "volume, vol", Value: math.NaN(), Usage: "音声効果:音量"},
		cli.Float64Flag{Name: "speed, spd", Value: math.NaN(), Usage: "音声効果:話速"},
		cli.Float64Flag{Name: "pitch, pit", Value: math.NaN(), Usage: "音声効果:高さ"},
		cli.Float64Flag{Name: "emphasis, emph", Value: math.NaN(), Usage: "音声効果:抑揚"},
		cli.BoolFlag{Name: "persist", Usage: "操作時の音声効果を残します"},
		cli.StringFlag{Name: "phrase, phrases", Usage: "フレーズ辞書相当の文字列を渡します。スペース区切りで複数のフレーズを フレーズ=発音 の羅列で指定します"},
		cli.BoolFlag{Name: "debug"},
	}
	app.Commands = []cli.Command{
		{
			Name: "test",
			Action: func(c *cli.Context) error {

				rog.EnableDebug()

				exe := `C:\Program Files (x86)\AHS\VOICEROID+\KiritanEX\VOICEROID.exe`
				kiri := vokiri.New(exe, "VOICEROID＋ 東北きりたん EX")

				if !kiri.IsRunning() {
					if err := kiri.Run(); err != nil {
						return fmt.Errorf("起動に失敗しました: %v\n", err)
					}
				}
				kiri.WaitForRun()

				kiri.Test()

				return nil
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		exe := c.String("exe")
		title := c.String("title")
		klose := c.Bool("close")
		record := c.String("record")
		recordOnce := c.String("record-once")
		volume := c.Float64("volume")
		speed := c.Float64("speed")
		pitch := c.Float64("pitch")
		emphasis := c.Float64("emphasis")
		persist := c.Bool("persist")
		phrase := c.String("phrase")

		text := strings.Join(c.Args(), " ")

		if c.Bool("debug") {
			rog.EnableDebug()
		}

		if len(exe) == 0 {
			exe = `C:\Program Files (x86)\AHS\VOICEROID+\KiritanEX\VOICEROID.exe`
			if _, err := os.Stat(exe); err != nil {
				// 32bit?
				exe = `C:\Program Files\AHS\VOICEROID+\KiritanEX\VOICEROID.exe`
				if _, err := os.Stat(exe); err != nil {
					fmt.Fprintf(os.Stderr, "VOICEROIDの実行ファイルが見つかりません。\n")
					fmt.Fprintf(os.Stderr, "フラグ --exe=\"～\" として実行ファイルを指定してください。\n")
					fmt.Fprintf(os.Stderr, "また、東北きりたん以外のVOICEROIDを使う場合は --title=\"～\" を指定してください。\n")
					return nil
				}
			}
		}

		if !termutil.Isatty(os.Stdin.Fd()) {
			data, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			decoded, err := ioutil.ReadAll(transform.NewReader(bytes.NewBuffer(data), japanese.ShiftJIS.NewDecoder()))
			if err == nil {
				text = string(decoded)
			} else {
				text = string(data)
			}
		}

		return run(exe, title, klose, record, recordOnce, volume, speed, pitch, emphasis, persist, text, phrase)
	}
	app.Run(os.Args)
	return
}

func run(exe, title string, klose bool, record string, recordOnce string, volume, speed, pitch, emphasis float64, persist bool, text, phrase string) error {
	var ips []string
	if extracteds, aftertext, err := vokiri.ExtractInstantPhrases(text); err != nil {
		return fmt.Errorf("インスタントフレーズの抽出に失敗しました(%q): %v", text, err)
	} else {
		ips = extracteds
		text = aftertext
	}

	if len(recordOnce) > 0 {
		if _, err := os.Stat(recordOnce); err == nil {
			skip := true
			txt := strings.Replace(strings.ToUpper(recordOnce), ".WAV", ".TXT", -1)
			if _, err := os.Stat(txt); err == nil {
				skip = !isContentChanged(text, txt)
			}
			if skip {
				return nil
			}
		}

		record = recordOnce
	}

	kiri := vokiri.New(exe, title)

	if !kiri.IsRunning() {
		if err := kiri.Run(); err != nil {
			return fmt.Errorf("起動に失敗しました: %v", err)
		}
	}
	kiri.WaitForRun()

	// register phrases
	{
		phrase = strings.Trim(strings.Replace(phrase, "　", " ", -1), " ")

		var pp []string
		if len(phrase) != 0 {
			phrase = strings.Replace(phrase, "＝", "=", -1)

			// check grammer of phrases
			pp = strings.Split(phrase, " ")
			for _, p := range pp {
				if p == "" {
					continue
				} else if !strings.Contains(p, "=") {
					return fmt.Errorf("オプション phrase(phrases) は = 区切りで指定します。(%q)", p)
				}
			}
		}

		if len(ips) != 0 {
			pp = append(pp, ips...)
		}

		if len(pp) != 0 {
			rog.Debug("フレーズ (persist=%v)", persist)
			for _, p := range pp {
				rog.Debug(p)
			}

			if err := kiri.ReloadPhraseDict(pp, persist); err != nil {
				return fmt.Errorf("フレーズ辞書編集の実行に失敗しました: %v\n", err)
			}
		}
	}

	var orgVolume, orgSpeed, orgPitch, orgEmphasis float64
	changed := false
	if !math.IsNaN(volume) {
		orgVolume = kiri.GetVolume()
		if err := kiri.SetVolume(volume); err != nil {
			return fmt.Errorf("音声効果:音量の設定に失敗しました: %v", err)
		}
		changed = true
	}
	if !math.IsNaN(speed) {
		orgSpeed = kiri.GetSpeed()
		if err := kiri.SetSpeed(speed); err != nil {
			return fmt.Errorf("音声効果:話速の設定に失敗しました: %v", err)
		}
		changed = true
	}
	if !math.IsNaN(pitch) {
		orgPitch = kiri.GetPitch()
		if err := kiri.SetPitch(pitch); err != nil {
			return fmt.Errorf("音声効果:高さの設定に失敗しました: %v", err)
		}
		changed = true
	}
	if !math.IsNaN(emphasis) {
		orgEmphasis = kiri.GetEmphasis()
		if err := kiri.SetEmphasis(emphasis); err != nil {
			return fmt.Errorf("音声効果:抑揚の設定に失敗しました: %v", err)
		}
		changed = true
	}
	if changed {
		time.Sleep(500 * time.Millisecond)
	}

	if record != "" {
		if err := kiri.Record(text, record); err != nil {
			return fmt.Errorf("音声ファイルの保存に失敗しました: %v", err)
		}
	} else {
		if err := kiri.Speak(text); err != nil {
			return fmt.Errorf("発話に失敗しました: %v", err)
		}
	}

	if !persist {
		if !math.IsNaN(volume) {
			kiri.SetVolume(orgVolume)
		}
		if !math.IsNaN(speed) {
			kiri.SetSpeed(orgSpeed)
		}
		if !math.IsNaN(pitch) {
			kiri.SetPitch(orgPitch)
		}
		if !math.IsNaN(emphasis) {
			kiri.SetEmphasis(orgEmphasis)
		}
	}

	if klose {
		kiri.Close()
	}

	return nil
}

func isContentChanged(text, filename string) bool {
	f, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()

	ret, err := ioutil.ReadAll(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))
	if err == nil {
		return string(ret) != text
	}
	return false

}
