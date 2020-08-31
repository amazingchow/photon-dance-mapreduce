PROJECT      := github.com/amazingchow/photon-dance-mapreduce
SRC          := $(shell find . -type f -name '*.go' -not -path "./vendor/*, ./apps/*")
MASTER       := mapreduce-master-service
WORKER       := mapreduce-worker-service
ALL_TARGETS  := $(MASTER) $(WORKER)
APP          := $(shell pwd)/apps/wc.so

LDFLAGS += -X "$(PROJECT)/utils.Plugin=$(APP)"

all: build

build: clean $(ALL_TARGETS)

$(MASTER): $(SRC)
	go build $(GOMODULEPATH)/$(PROJECT)/cmd/$@

$(WORKER): $(SRC)
	go build -ldflags '$(LDFLAGS)' $(GOMODULEPATH)/$(PROJECT)/cmd/$@

lint:
	@golangci-lint run --skip-dirs=api --deadline=5m

pb-fmt:
	@clang-format -i ./pb/*.proto

test:
	go test -count=1 -v -p 1 $(shell go list ./backend/storage...)
	go test -count=1 -v -p 1 $(shell go list ./master...)
	go test -count=1 -v -p 1 $(shell go list ./utils...)
	go test -count=1 -v -p 1 $(shell go list ./worker...)

clean:
	rm -f $(ALL_TARGETS)

.PHONY: all build clean
