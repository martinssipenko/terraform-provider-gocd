.DEFAULT_GOAL := gocd

CWD:=$(CURDIR)
TRAVIS_TAG ?= dev

fmt: format
	git add .
	git commit -m "Autocommit for 'gofmt -w'."

format:
	gofmt -w -s .

gocd.%: init
	(cd $(CWD)/cli && \
		mkdir -p $(CWD)/build/$*/ && \
		GOOS=$* go build -o $(CWD)/build/$*/gocd && \
		cd $(CWD)/build/$* && \
		tar -czvf $(CWD)/target/gocd-$*-$(TRAVIS_TAG).tgz ./gocd)

init:
	mkdir -p $(CWD)/target/
	mkdir -p $(CWD)/build/

clean:
	rm -rf $(CWD)/target
	rm -rf $(CWD)/build