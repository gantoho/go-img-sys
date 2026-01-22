@echo off
REM Build and Run script for Windows

setlocal enabledelayedexpansion

if "%1"=="" (
    echo Usage: .\build.bat [command]
    echo.
    echo Commands:
    echo   build    - Build the application
    echo   run      - Build and run the application
    echo   clean    - Clean build artifacts
    echo   help     - Show this help message
    exit /b 0
)

if "%1"=="build" (
    echo Building application...
    go build -o image-sys.exe main.go
    if !errorlevel! equ 0 (
        echo Build successful: image-sys.exe
    ) else (
        echo Build failed!
        exit /b 1
    )
)

if "%1"=="run" (
    call :build
    if !errorlevel! equ 0 (
        echo Starting application on http://localhost:3128
        image-sys.exe
    )
)

if "%1"=="clean" (
    echo Cleaning build artifacts...
    del /f /q image-sys.exe 2>nul
    go clean
    echo Clean complete
)

if "%1"=="help" (
    echo Usage: .\build.bat [command]
    echo.
    echo Commands:
    echo   build    - Build the application
    echo   run      - Build and run the application
    echo   clean    - Clean build artifacts
    echo   help     - Show this help message
)

goto :eof

:build
go build -o image-sys.exe main.go
exit /b !errorlevel!
