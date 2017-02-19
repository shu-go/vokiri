set GOOS=windows
set GOARCH=386
pushd cmd\vokiri 
go build -ldflags "-s -w" %*
del vokiri.zip 2>nul
zip vokiri.zip vokiri.exe README.txt LICENSE.txt sample.bat launch.bat
popd
move cmd\vokiri\vokiri.zip .
