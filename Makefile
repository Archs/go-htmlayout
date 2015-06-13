all:install

zsyscall_htmlayout.go:htmlayout.go element.go value.go
	mksyscall_dll -p gohl $^ > $@

install:htmlayout.go zsyscall_htmlayout.go
	go install -x .

clean:
	rm -rf zsyscall_*.go