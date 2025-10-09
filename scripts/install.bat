@echo off
REM Installation script for my-context (Windows cmd.exe)
REM Installs to %USERPROFILE%\bin without requiring admin privileges

setlocal enabledelayedexpansion

set "INSTALL_DIR=%USERPROFILE%\bin"
set "BINARY_NAME=my-context.exe"

echo Installing my-context for Windows...

REM Detect existing installation
if exist "%INSTALL_DIR%\%BINARY_NAME%" (
    echo Existing installation found.
    
    REM Get current version
    "%INSTALL_DIR%\%BINARY_NAME%" --version 2>nul
    
    REM Prompt for upgrade
    set /p UPGRADE="Existing installation found. Upgrade? (y/N): "
    if /i not "!UPGRADE!"=="y" (
        echo Installation cancelled.
        exit /b 0
    )
    
    REM Backup old binary
    copy "%INSTALL_DIR%\%BINARY_NAME%" "%INSTALL_DIR%\%BINARY_NAME%.bak" >nul
    echo Backed up existing binary
)

REM Create install directory
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

REM Copy binary (assume binary provided as %1 or from bin/)
if not "%~1"=="" (
    copy "%~1" "%INSTALL_DIR%\%BINARY_NAME%" >nul
) else (
    REM Auto-detect binary from bin/
    if exist "bin\my-context-windows-amd64.exe" (
        copy "bin\my-context-windows-amd64.exe" "%INSTALL_DIR%\%BINARY_NAME%" >nul
    ) else (
        echo Error: Binary not found at bin\my-context-windows-amd64.exe
        echo Run 'scripts\build-all.sh' first to build binaries
        exit /b 1
    )
)

REM Add to PATH if not already present
echo %PATH% | find /i "%INSTALL_DIR%" >nul
if errorlevel 1 (
    echo Adding %INSTALL_DIR% to user PATH...
    setx PATH "%PATH%;%INSTALL_DIR%" >nul
    echo Added to PATH. Restart terminal for changes to take effect.
) else (
    echo %INSTALL_DIR% already in PATH
)

REM Verify installation
"%INSTALL_DIR%\%BINARY_NAME%" --version >nul 2>&1
if errorlevel 1 (
    echo.
    echo Installation verification failed
    exit /b 1
)

echo.
echo Installation complete!
"%INSTALL_DIR%\%BINARY_NAME%" --version
echo.
echo Note: ~/.my-context/ data directory is preserved (separate from binary)
echo Restart your terminal for PATH changes to take effect

REM Remove backup if successful
if exist "%INSTALL_DIR%\%BINARY_NAME%.bak" del "%INSTALL_DIR%\%BINARY_NAME%.bak" >nul

endlocal

