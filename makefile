.PHONY: build build-win dev run clean

PROJECT="game-gorl"
BUILD_PATH="./build"
USE_PACKFILE=1

init:
	mkdir build
	mkdir runtime

build:
	mkdir -p $(BUILD_PATH)
	if [ $(USE_PACKFILE) -eq 1 ]; then \
		go run cmd/tool/main.go packer ./assets $(BUILD_PATH)/data.pack; \
	else \
		cp -r assets/* $(BUILD_PATH); \
	fi
	go build -o $(BUILD_PATH)/$(PROJECT) -v cmd/game/main.go

build-debug:
	mkdir -p $(BUILD_PATH)
	cp -r assets/* $(BUILD_PATH)
	CGO_CFLAGS='-O0 -g' go build -a -v -gcflags="all=-N -l" -o $(BUILD_PATH)/$(PROJECT) cmd/game/main.go 

run:
	cd $(BUILD_PATH); ./$(PROJECT)

dev:
	@make clean && make build && make run || echo "failure!"

clean:
	rm -r $(BUILD_PATH)/*

lint:
	nilaway main.go
