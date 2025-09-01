default: build

gobin:
	echo $(GOBIN)

build:
	go build -o jn .

install:
	go build -o jn . && go install . && rm jn

uninstall:
	rm -f $(GOBIN)/jn

tokei:
	rm TOKEI && tokei . > TOKEI
