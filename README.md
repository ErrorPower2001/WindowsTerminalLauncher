# Description

A Windows Terminal profiles launcher.

# Download, Install, Setup

To download the prebuild executable, please read [Releases](https://github.com/ErrorPower2001/WindowsTerminalLauncher/releases/).

After download, move the executable to any directory and add it to the Path environment variable.

Next, add the executable to Windows Terminal as a new profile and set it as the default option.

If you use Scoop Installer, you can run:

```
scoop install https://raw.githubusercontent.com/ErrorPower2001/WindowsTerminalLauncher/master/WindowsTerminalLauncher.json
```



# Build

```
go build -o wtlauncher.exe .\main.go
```
