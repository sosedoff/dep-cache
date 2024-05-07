build:
	go build

publish:
	goreleaser --skip-validate --rm-dist
