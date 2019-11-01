build:
	go build
	mv luxafor-macos Luxafor.app/Contents/MacOS/

run: build
	./Luxafor.app/Contents/MacOS/luxafor-macos