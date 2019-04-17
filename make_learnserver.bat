@echo off

set VERSION=v.0.1.0
echo version=%VERSION%

set GOPATH=D:\teaching
set GOPATH_BAK=%GOPATH%
set GOPATH=%GOPATH%;%~dp0;%~dp0\src;

set GOOS=windows
set GOARCH=amd64

echo start install sever ...
go install -gcflags " -N -l" LearnServer

echo make ok

set GOPATH=%GOPATH_BAK%

pause