.DEFAULT_GOAL := gocd

SHELL:=/bin/bash
CWD:=$(CURDIR)
TRAVIS_TAG ?= dev

fmt:
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

test: vet lint
	bash ../scripts/go-test.sh

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

lint:
	golint .
