@echo off
REM Setup script for RestySched (Windows)

echo === RestySched Setup ===
echo.

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed. Please install Go 1.23 or higher.
    exit /b 1
)

for /f "tokens=*" %%i in ('go version') do set GO_VERSION=%%i
echo [OK] Go is installed: %GO_VERSION%
echo.

REM Install Templ CLI
echo Installing Templ CLI...
go install github.com/a-h/templ/cmd/templ@latest
echo [OK] Templ CLI installed
echo.

REM Download dependencies
echo Downloading Go dependencies...
go mod download
echo [OK] Dependencies downloaded
echo.

REM Generate templates
echo Generating Templ templates...
templ generate
echo [OK] Templates generated
echo.

REM Create .env file if it doesn't exist
if not exist .env (
    echo Creating .env file from example...
    copy .env.example .env
    echo [OK] .env file created
    echo.
    echo [WARNING] Please update .env file with your n8n webhook URL!
) else (
    echo [OK] .env file already exists
)

echo.
echo === Setup Complete! ===
echo.
echo Next steps:
echo 1. Update .env file with your n8n webhook URL
echo 2. Run: go run cmd/server/main.go
echo 3. Access the app at http://localhost:8080
echo.
pause
