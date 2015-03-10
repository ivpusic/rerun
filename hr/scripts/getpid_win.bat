@ECHO OFF

call netstat -aon | c:\windows\system32\find.exe /i "listening" | c:\windows\system32\find.exe "%1" | gawk '{ print $5;exit }'

exit 0