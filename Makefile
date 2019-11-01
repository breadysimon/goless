.PHONY: test


default: test

help:
	@echo 'Usage:'
	@echo '    make test            Run tests on a compiled project.'
	@echo

test:
	go test ./...

