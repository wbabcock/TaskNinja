APP_NAME = taskninja
MAC_NAME = $(APP_NAME)_macos-aarch64
WIN_NAME = $(APP_NAME)_win-amd64.exe
MAC_TARGET = bin/$(MAC_NAME)
WIN_TARGET = bin/$(WIN_NAME)

# ./$(TARGET) add this is the first task project:test +NEW +MINE due:eow

run: clean default
	./$(MAC_TARGET)

default: $(MAC_TARGET)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(WIN_TARGET) ./src

clean:
	rm -f bin/*

install: clean default
	sudo cp ./bin/taskninja_macos-aarch64 /usr/local/bin/taskninja

build-win: clean default
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(WIN_TARGET) ./src

$(MAC_TARGET):
	go build -o $(MAC_TARGET) ./src
