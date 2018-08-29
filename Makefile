TEST_PACKAGE = test/...
release:
	git tag -a v$(VER) -m 'v$(VER)'
	git push origin --tags
	echo Tagged release with $(VER)
.PHONY: release

test:
	go test ./$(TEST_PACKAGE) -v -coverprofile=coverage.out
.PHONY: test