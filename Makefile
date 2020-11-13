BIN=a320sim a320cli a320ecam

.PHONY: all
all: build

.PHONY: build
build:
	$(foreach bin,$(BIN),go build ./cmd/$(bin);)

.PHONY: clean
	rm -f $(BIN)
