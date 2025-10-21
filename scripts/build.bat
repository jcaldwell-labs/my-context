@echo off
REM build.bat - Windows cmd.exe build script with proper metadata
REM Usage: scripts\build.bat

echo Building my-context for Windows...

REM Get version from git tag
for /f %%i in ('git describe --tags --exact-match 2^>nul') do set VERSION=%%i
if "%VERSION%"=="" (
    for /f %%i in ('git describe --tags --always') do set VERSION=%%i
)

REM Get build time (Windows datetime format)
for /f "tokens=1-4 delims=/ " %%a in ('date /t') do set BUILD_DATE=%%c-%%a-%%b
for /f "tokens=1-2 delims=: " %%a in ('time /t') do set BUILD_TIME=%%a:%%b
set BUILD_TIME=%BUILD_DATE%T%BUILD_TIME%:00Z

REM Get git commit
for /f %%i in ('git rev-parse --short HEAD') do set GIT_COMMIT=%%i

echo    Version: %VERSION%
echo    Build: %BUILD_TIME%
echo    Commit: %GIT_COMMIT%
echo.

REM Build with ldflags (note the = sign and quotes!)
go build -ldflags="-X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME% -X main.GitCommit=%GIT_COMMIT%" -o my-context.exe cmd/my-context/main.go

if %ERRORLEVEL% neq 0 (
    echo Build failed!
    exit /b 1
)

echo.
echo Build complete: my-context.exe
echo.
my-context.exe --version
echo.
echo To use: copy my-context.exe to your PATH or run .\my-context.exe
