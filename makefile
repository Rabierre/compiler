CC = go

test:
	$(CC) test ./...
clean:
	$(CC) clean
ifneq ("$(wildcard *.orig)","")
	rm -f *.orig
endif
ifneq ("$(wildcard *.a)","")
	rm -f *.a
endif

build:
	$(CC) build -o compiler .
install:
	$(CC) get github.com/rabierre/compiler
.PONEY: clean build test install
