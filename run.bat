go build -o go-websockets.exe ./cmd/web/. || exit /b
go-websockets.exe
-dbname=go-chatroom
-dbuser=postgres
-dbpass=andrew00
-cache=false
-production=false