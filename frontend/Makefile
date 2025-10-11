# windows 下装个 https://www.msys2.org/
# 然后安装 pacman -S make
# 然后设置下PATH环境变量，C:\msys64\usr\bin，windows下就可以用make了

# 协议生成
# https://buf.build/docs/cli/installation/
# GOBIN=/usr/local/bin go install github.com/bufbuild/buf/cmd/buf@v1.54.0

ifneq (,$(filter $(OS),Windows_NT MINGW64))
EXE = .exe
MKDIR = mkdir
else
MKDIR = mkdir -p
endif

RM = rm -rf
	
all: build

dev-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.8
	go install github.com/favadi/protoc-go-inject-tag@latest
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

proto:
	protoc \
		--plugin=./node_modules/.bin/protoc-gen-ts_proto \
		--proto_path=./proto \
		--ts_proto_out=./src/lib/proto/ \
		--ts_proto_opt="esModuleInterop=true,forceLong=long" \
		./proto/**/*.proto	

tidy:
	go mod tidy -C .

config: | build-exportcsv
	$(BUILD_OUTPUT)/exportcsv$(EXE) "export"	

build:
	npm: build - ,

.PHONY: proto dev-tools build