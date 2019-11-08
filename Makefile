build:
	go build
	mv luxafor-macos Luxafor.app/Contents/MacOS/
	./update-icon.sh

run: build
	./Luxafor.app/Contents/MacOS/luxafor-macos