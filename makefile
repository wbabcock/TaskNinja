TARGET = bin/taskninja_macos-aarch64
WIN_TARGET = bin/taskninja_win-amd64.exe

# ./$(TARGET) add this is the first task project:test +NEW +MINE due:eow

run: clean default
	./$(TARGET)

default: $(TARGET)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(WIN_TARGET) ./src

clean:
	rm -f bin/*

install: clean default
	sudo cp ./bin/taskninja_macos-aarm /usr/local/bin/taskninja

build-win: clean default
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(WIN_TARGET) ./src

$(TARGET):
	go build -o $(TARGET) ./src
