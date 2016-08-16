start cmd /K crosssvr.exe
start cmd /K accountsvr.exe
start cmd /K chatsvr.exe
start cmd /K gamesvr.exe
choice /t 1 /d y /n >nul
start cmd /K battlesvr.exe -port 8101

