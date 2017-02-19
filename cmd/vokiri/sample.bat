@echo off
setlocal ENABLEEXTENSIONS
rem %~dp0 までが、このファイルがあるフォルダ＋￥
set DEST=%~dp0kiritan_
rem set DEBUG=--debug

pushd %~dp0

if "%1"=="help" (
    call :help
) else if "%1"=="gen" (
    call :gen
) else if "%1"=="clean" (
    call :clean
) else  (
    call :sample
)

popd
endlocal
echo on

@exit /b 0


:sample

vokiri --volume=1 --speed=0.95 --pitch=1 --emph=1.2 --persist

vokiri %DEBUG% 今日は%date:~0,4%年%date:~5,2%月%date:~8,2%日です。
vokiri %DEBUG% VOKIRIに興味を持っていただきありがとうございます。
vokiri %DEBUG% VOKIRIは、コマンドラインからVOICEROIDを操作するためのソフトウェアです。
vokiri %DEBUG% 手作業を削減したり、ほかのソフトからの橋渡し役になったりすることを期待しています。

vokiri %DEBUG% さて、VOKIRIでは、若干むりやりにVOICEROIDのウィンドウを操作していることもあり＜＜だって、それしか方法がないんですもん｜＞＞
vokiri %DEBUG% 稀に動きが止まってしまうことがあります。
vokiri %DEBUG% その場合は、もう一回同じことを実行してみてください。多分動くようになると思います。
vokiri %DEBUG% それでもだめなら、コマンドラインのオプションとして、＜＜--debug｜はいふんはいふんでばっぐ＞＞を指定して、止まった個所を作者に教えてあげてください。

vokiri %DEBUG% あ、あと、東北きりたん以外のVOICEROIDでの動作状況についても、教えていただけると幸いです。
vokiri %DEBUG% ＜＜--exe｜はいふんはいふん いーえっくすいー＞＞、で、動作させたいVOICEROIDの実行ファイルを、
vokiri %DEBUG% ＜＜--title｜はいふんはいふん たいとる＞＞、で、そのVOICEROIDのウィンドウタイトル（"VOICEROID 琴葉茜"等）を指定してみてください。

vokiri %DEBUG% ではでは。

ping -w 0 -n 2 0.0.0.0>nul
vokiri %DEBUG% --volume=0.5 --pitch=2 --emphasis=2 ｷﾘﾀﾝｶﾜｲｲﾔｯﾀｰ
vokiri %DEBUG% --volume=0.5 --pitch=2 --emphasis=2 ｾﾔﾅｰ　

exit /b 0


:clean

del kiritan*

exit /b 0


:gen

echo %TIME%

rem
rem 紹介動画向けの音声ファイル作成。
rem

vokiri --volume=1 --speed=0.95 --pitch=1 --emph=1.2 --persist

rem
rem 挨拶
rem
vokiri %DEBUG% --record-once="%DEST%0001.wav" "こんにちは、（（東北きりたん｜(spd:0.8)トーホク(pit:1.20)キ!リタン$1_1））（（です｜(pit:1.2)デ｜スD））。 "
vokiri %DEBUG% --record-once="%DEST%0002.wav" この動画では、VOKIRIというVOICEROID支援ソフトの紹介と、導入までの簡単な説明をします。
vokiri %DEBUG% --record-once="%DEST%0003.wav" VOKIRIは、VOICEROIDをお使いの方向けのソフトで、上に挙げる機能を持っています。
vokiri %DEBUG% --record-once="%DEST%0004.wav" ちなみに、現在は東北きりたんの操作画面に合わせて作りこんでいますが、将来的には他のVOICEROIDにも対応する…かもしれません。
rem コマンドラインでの読み上げ
vokiri %DEBUG% --record-once="%DEST%0005.wav" 機能そのいち。コマンドラインでVOICEROID＋ 東北きりたんに読み上げさせる
vokiri %DEBUG% --record-once="%DEST%0006.wav" 基本機能です。
rem イントネーション等の調音
vokiri %DEBUG% --record-once="%DEST%0007.wav" 機能そのに。イントネーション等の調音をコマンドラインで指定する
vokiri %DEBUG% --record-once="%DEST%0008.wav" 通常のコマンドラインツール等ではあまり実現していない、コマンドラインからの調音を行えます。
rem 音声ファイルへの保存
vokiri %DEBUG% --record-once="%DEST%0009.wav" 機能そのさん。読み上げた内容を音声ファイルとして保存する
vokiri %DEBUG% --record-once="%DEST%0010.wav" 多くのツールで実現できていることですが、当然VOKIRIでも行えます。
vokiri %DEBUG% --record-once="%DEST%0011.wav" この動画も、この方法で ＜＜東北きりたん｜わたし＞＞に読み上げさせています。
rem 無変更の音声ファイル保存のスキップ
vokiri %DEBUG% --record-once="%DEST%0012.wav" 機能そのよん。無変更の音声ファイル保存のスキップ
vokiri %DEBUG% --record-once="%DEST%0013.wav" 複数のコマンドラインで音声ファイルを一気に書き出したい場合、出力には時間がかかります。
vokiri %DEBUG% --record-once="%DEST%0014.wav" この機能を使うと、いったん音声ファイルを出力したうえで一部を変更した場合の出力スピードを大幅に改善します。
rem
vokiri %DEBUG% --record-once="%DEST%0015.wav" VOKIRIはこういった機能を持っています。
vokiri %DEBUG% --record-once="%DEST%0016.wav" "（（活用してくださいね｜(pit:0.95)カツヨウ シテ (pit:1.1)ク^ダサ!イ(spd:2)(emph:0.5)ネ^ー<R>））。"

rem
rem ダウンロード
rem
vokiri %DEBUG% --record-once="%DEST%0100.wav" ダウンロードのしーかーたー
vokiri %DEBUG% --record-once="%DEST%0101.wav" まずは、お手元に、わたくし東北きりたんをご用意ください。
vokiri %DEBUG% --record-once="%DEST%0102.wav" --emph=2 もし、万が一無いようでしたら、これを機にご購入ください！
vokiri %DEBUG% --record-once="%DEST%0103.wav" --speed=0.8 こほん、
vokiri %DEBUG% --record-once="%DEST%0104.wav" そして、このVOKIRIをダウンロードしてください。
vokiri %DEBUG% --record-once="%DEST%0105.wav" ダウンロード場所は、動画のコメント欄に記載してあります。
vokiri %DEBUG% --record-once="%DEST%0106.wav" ダウンロードしたら、ZIPファイルを展開します。

rem
rem 使ってみる
rem
vokiri %DEBUG% --record-once="%DEST%0200.wav" 使ってみましょう！
vokiri %DEBUG% --record-once="%DEST%0201.wav" ダウンロードができたら、中にある＜＜sample｜サンプル＞＞というファイルを実行してみてください。
vokiri %DEBUG% --record-once="%DEST%0202.wav" 私、東北きりたんが何やら読み上げはじめます。
vokiri %DEBUG% --record-once="%DEST%0203.wav" このファイルには、VOKIRIを使って私に読み上げさせる、という処理が記述されています。
vokiri %DEBUG% --record-once="%DEST%0204.wav" そのほかにも、この動画に使うための音声ファイル出力のコマンドが詰まっています。
vokiri %DEBUG% --record-once="%DEST%0205.wav" 同様に＜＜同梱｜どうこん＞＞されている ＜＜README｜りーどみー＞＞ ともども、参考になるかと思います。

rem
rem さいごに
rem
vokiri %DEBUG% --record-once="%DEST%0300.wav" こんなところです！
vokiri %DEBUG% --record-once="%DEST%0301.wav" これでVOKIRIの紹介と説明を終わります。
vokiri %DEBUG% --record-once="%DEST%0302.wav" ごくごく簡単な紹介になりましたが、ご容赦ください。
vokiri %DEBUG% --record-once="%DEST%0303.wav" ＜＜同梱｜どうこん＞＞されている ＜＜README｜りーどみー＞＞ には、より詳細な機能の使い方が記載されています。
vokiri %DEBUG% --record-once="%DEST%0304.wav" ほかのコマンドが出力したテキスト(標準出力)を受け取って読み上げさせる手段も紹介しています。
vokiri %DEBUG% --record-once="%DEST%0305.wav" 活用していただけると嬉しいです。
vokiri %DEBUG% --record-once="%DEST%0306.wav" ではではー！

rem
rem おまけ
rem
vokiri --record-once="%DEST%omake_kwaii.wav" --pitch=2 --emphasis=2 ｷﾘﾀﾝｶﾜｲｲﾔｯﾀｰ　
vokiri --record-once="%DEST%omake_seya.wav" --pitch=2 --emphasis=2 ｾﾔﾅｰ　

vokiri おしまーーい！

dir /b *.wav > kiritan.m3u

echo %TIME%

exit /b 0


:help

echo HELP:
echo   %~n0 gen: 動画用のWAVファイルを出力します。テキストファイルを出力する設定の場合、変更したものだけを処理します。
echo   %~n0 clean: 出力されたファイル（kiritan～）を削除します。

exit /b 0

