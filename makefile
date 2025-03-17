build:
	echo "building web server application"
	go build -o bin/codex github.com/gary-norman/forum/cmd/web

run:
	echo "running web server application"
	bin/codex
