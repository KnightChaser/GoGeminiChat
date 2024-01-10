# one stop build and execution
$scriptPath = $PSScriptRoot
Set-Location $scriptPath
go build "./main.go"
./main.exe