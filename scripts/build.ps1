# Build From Source
$loc = "$HOME\AppData\Local\doko"

go run scripts/date.go >> date.txt

$LATEST_VERSION=git describe --abbrev=0 --tags
$DATE=cat date.txt

# Build
go mod tidy
go build -o doko.exe -ldflags "-X main.version=$LATEST_VERSION -X main.versionDate=$DATE"

# Setup
$BIN = "$loc\bin"
New-Item -ItemType "directory" -Path $BIN
Move-Item doko.exe -Destination $BIN
[System.Environment]::SetEnvironmentVariable("Path", $Env:Path + ";$BIN", [System.EnvironmentVariableTarget]::User)

if (Test-Path -path $loc) {
    Write-Host "Doko was built successfully, refresh your powershell and then run 'doko --help'" -ForegroundColor DarkGreen
} else {
    Write-Host "Build failed" -ForegroundColor Red
}
