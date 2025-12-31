@echo off
REM Build script for GoStructure App (Windows)

SET APP_NAME=app.exe
SET BUILD_DIR=.\build\bin
SET MAIN_PATH=.\cmd\app

REM Create build directory
if not exist %BUILD_DIR% mkdir %BUILD_DIR%

echo Building %APP_NAME%...

REM Build for Windows
go build -o "%BUILD_DIR%\%APP_NAME%" %MAIN_PATH%

IF %ERRORLEVEL% EQU 0 (
    echo Build complete: %BUILD_DIR%\%APP_NAME%
) ELSE (
    echo Build failed!
    exit /b 1
)
