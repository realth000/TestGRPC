@echo off
echo build server...
go build --buildmode=exe -o client ./greeter_client || goto fail
echo build client...
go build --buildmode=exe -o server ./greeter_server || goto fail
echo done

pause
exit /b 0


:fail
echo build failed
pause
