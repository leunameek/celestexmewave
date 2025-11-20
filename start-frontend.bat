@echo off
REM CelestexMewave Frontend Start Script (Windows)
REM This script starts a local development server for the frontend

setlocal enabledelayedexpansion

echo.
echo 🚀 Starting CelestexMewave Frontend...
echo.

REM Check if Python 3 is installed
python --version >nul 2>&1
if errorlevel 0 (
    echo ✓ Python found
    echo.
    echo 🔧 Starting frontend server on http://localhost:3000
    echo.
    python -m http.server 3000
    exit /b 0
)

REM Check if Node.js is installed
node --version >nul 2>&1
if errorlevel 0 (
    echo ✓ Node.js found
    echo.
    
    REM Check if http-server is installed globally
    http-server --version >nul 2>&1
    if errorlevel 0 (
        echo ✓ http-server found
        echo.
        echo 🔧 Starting frontend server on http://localhost:3000
        echo.
        http-server -p 3000
        exit /b 0
    ) else (
        echo 📦 Installing http-server...
        call npm install -g http-server
        echo.
        echo 🔧 Starting frontend server on http://localhost:3000
        echo.
        call http-server -p 3000
        exit /b 0
    )
)

REM If we get here, no suitable server was found
echo ❌ Error: No suitable server found!
echo.
echo Please install one of the following:
echo   • Python 3 (https://www.python.org/downloads/)
echo   • Node.js (https://nodejs.org/)
echo.
echo Then run this script again.
echo.
pause
exit /b 1
