
TARGETS = all test clean

.PHONY: $(TARGETS) fmt

$(TARGETS):
	$(MAKE) -C cmd/ $@

fmt:
	gofmt -w .
