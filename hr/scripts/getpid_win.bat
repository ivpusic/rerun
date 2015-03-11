@ECHO OFF

setlocal
set sys32dir=%SystemRoot%\system32

for /F "tokens=5 delims= " %%i in ('%sys32dir%\netstat.exe -aon ^| %sys32dir%\find.exe /i "LISTENING" ^| %sys32dir%\find.exe "%1"') do (
  echo %%i
  goto end
)
:end
exit /b 0
