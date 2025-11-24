@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul

REM Qiniu Uploader Windows Build Script

echo Starting Qiniu Uploader build...

REM Check if Go is installed
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Go compiler not found
    echo Please install Go 1.21 or higher
    echo Download: https://golang.org/dl/
    exit /b 1
)

REM Check Go version
for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
set GO_VERSION=%GO_VERSION:go=%

REM Parse version number
for /f "tokens=1,2 delims=." %%a in ("%GO_VERSION%") do (
    set MAJOR=%%a
    set MINOR=%%b
)

REM Check if meets minimum version requirement (Go 1.21)
if %MAJOR% lss 1 (
    echo ERROR: Go version too low
    echo Current version: %GO_VERSION%
    echo Required version: 1.21 or higher
    exit /b 1
) else if %MAJOR% equ 1 if %MINOR% lss 21 (
    echo ERROR: Go version too low
    echo Current version: %GO_VERSION%
    echo Required version: 1.21 or higher
    exit /b 1
)

echo Go version check passed: %GO_VERSION%

REM Download dependencies
echo Downloading dependencies...
go mod download
if %errorlevel% neq 0 (
    echo ERROR: Dependency download failed
    exit /b 1
)

REM Build program
echo Building program...
go build -o qu.exe .\cmd\qiniu-uploader
if %errorlevel% neq 0 (
    echo ERROR: Build failed
    exit /b 1
)

echo Build successful!
echo.
echo Usage:
echo   1. Initialize config: qu.exe config init
echo   2. Upload file: qu.exe upload
echo   3. View help: qu.exe --help
echo.
echo Tip: You can move qu.exe to any directory in your PATH for easy access
echo.
echo Qiniu Uploader is ready!

endlocal