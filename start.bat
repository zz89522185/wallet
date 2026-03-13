@echo off
chcp 65001 >nul
set "PROJECT_DIR=%~dp0"
set "BIN_DIR=%PROJECT_DIR%bin"

:: Locate go.exe
for /f "delims=" %%i in ('where go 2^>nul') do set "GO_EXE=%%i" & goto :found_go
echo ERROR: go.exe not found in PATH
pause
exit /b 1
:found_go

:: Force Windows/amd64 build target
set "GOOS=windows"
set "GOARCH=amd64"

echo ========== Building wallet services ==========

if not exist "%BIN_DIR%" mkdir "%BIN_DIR%"

echo [1/2] Building RPC service...
"%GO_EXE%" build -o "%BIN_DIR%\wallet-rpc.exe" "%PROJECT_DIR%service\wallet\rpc\wallet.go"
if %errorlevel% neq 0 (
    echo ERROR: RPC build failed
    pause
    exit /b 1
)

echo [2/2] Building API service...
"%GO_EXE%" build -o "%BIN_DIR%\wallet-api.exe" "%PROJECT_DIR%service\wallet\api\wallet.go" "%PROJECT_DIR%service\wallet\api\doc.go"
if %errorlevel% neq 0 (
    echo ERROR: API build failed
    pause
    exit /b 1
)

echo Build complete.
echo.
echo ========== Starting wallet services ==========

echo [1/2] Starting RPC service (port 8080)...
start "wallet-rpc" cmd /k "cd /d "%PROJECT_DIR%" && "%BIN_DIR%\wallet-rpc.exe" -f service\wallet\rpc\etc\wallet.yaml"

timeout /t 3 /nobreak >nul

echo [2/2] Starting API service (port 8888)...
start "wallet-api" cmd /k "cd /d "%PROJECT_DIR%" && "%BIN_DIR%\wallet-api.exe" -f service\wallet\api\etc\wallet-api.yaml"

echo.
echo All services started.
echo   RPC: 0.0.0.0:8080
echo   API: http://localhost:8888
echo   Swagger: http://localhost:8888/swagger/index.html
echo.
echo To stop all services, run: stop.bat
