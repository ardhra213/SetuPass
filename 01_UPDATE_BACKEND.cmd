@echo off
cd /d "%~dp0"
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%~dp0scripts\Update-Backend.ps1"
echo.
if errorlevel 1 echo STOP: Keep this output and share the error.
pause
