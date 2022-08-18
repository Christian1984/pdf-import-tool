@ECHO OFF
SET noguiflag=-H=windowsgui
IF %1.==dev. SET noguiflag=

ECHO create winres meta data...
go-winres make

ECHO build project...
go build -o pdf-importer.exe -ldflags "-s -w %noguiflag%" -v .