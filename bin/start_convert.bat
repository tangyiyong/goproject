..\..\tools\xls2csv\xls2csv.exe ..\..\cehua\ConfigData .\csvtemp ..\..\tools\xls2csv\heros.x2c
for /r .\csvtemp %%i in (*.csv) do ..\..\tools\xls2csv\iconv.exe -f GBK -t UTF-8 %%i > .\csv\%%~nxi
rmdir /s /q csvtemp
