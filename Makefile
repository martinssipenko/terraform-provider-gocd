.PHONY := format

fmt: format
	git add .
	git commit -m "Autocommit for 'gofmt -w'."

format:
	gofmt -w -s .
	$(MAKE) -C ./gocd/ format
