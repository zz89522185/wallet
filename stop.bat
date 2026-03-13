@echo off
chcp 65001 >nul
echo Stopping wallet services...

:: Kill wallet-rpc.exe
tasklist /FI "IMAGENAME eq wallet-rpc.exe" 2>nul | findstr /I "wallet-rpc.exe" >nul
if %errorlevel%==0 (
    echo Killing wallet-rpc.exe...
    taskkill /F /IM wallet-rpc.exe >nul 2>&1
) else (
    echo wallet-rpc.exe not running.
)

:: Kill wallet-api.exe
tasklist /FI "IMAGENAME eq wallet-api.exe" 2>nul | findstr /I "wallet-api.exe" >nul
if %errorlevel%==0 (
    echo Killing wallet-api.exe...
    taskkill /F /IM wallet-api.exe >nul 2>&1
) else (
    echo wallet-api.exe not running.
)

echo All services stopped.
