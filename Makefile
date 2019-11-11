build:
	CGO_CFLAGS=-mmacosx-version-min=10.12 CGO_LDFLAGS=-mmacosx-version-min=10.12 go build -ldflags -v
	mv luxafor-macos Luxafor.app/Contents/MacOS/
	./update-icon.sh

run: build
	./Luxafor.app/Contents/MacOS/luxafor-macos

install: build
	cp -r ./Luxafor.app /Applications