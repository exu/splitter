default:
	go build -o bin/amd64/splitter

install:
	go build -o bin/amd64/splitter
	cp bin/x64/splitter /usr/bin/
