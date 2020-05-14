PROJECT      := amazingchow/mapreduce
SRC          := $(shell find . -type f -name '*.go' -not -path "./vendor/*, ./apps/*")
PB_SRC       := $(shell find . -type f -name '*.proto' -not -path "./vendor/*")
TARGETS      := mapreduce-master-service mapreduce-worker-service
ALL_TARGETS  := $(TARGETS)

all: build

build: clean $(ALL_TARGETS)

$(TARGETS): $(SRC)
	go build $(GOMODULEPATH)/$(PROJECT)/cmd/$@

lint:
	@golangci-lint run --skip-dirs=api --deadline=5m

pb-fmt:
	@clang-format -i ./pb/*.proto

test:
	go test -count=1 -v -p 1 $(shell go list ./ch/...| grep -v /vendor/)

clean:
	rm -f $(ALL_TARGETS)

.PHONY: all build clean
