@ECHO OFF
SET noguiflag=-H=windowsgui
IF %1.==dev. SET noguiflag=

ECHO create winres meta data...
go-winres make

ECHO generate bindata
go generate -v .\res

ECHO get version string
git describe --tags>versionstr.txt
SET /p versionstr=<versionstr.txt
del versionstr.txt

ECHO build project...
go build -o pdf-importer.exe -ldflags "-s -w -X=main.BuildVersion=%versionstr% %noguiflag%" -v .