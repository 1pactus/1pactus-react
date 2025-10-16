ifneq (,$(filter $(OS),Windows_NT MINGW64))
EXE = .exe
MKDIR = mkdir
else
MKDIR = mkdir -p
endif

RM = rm -rf

BUILD_OUTPUT?=$(shell pwd)/build

APPS = onepacd
YAML_CFGS = onepacd
	
all: tidy build cp-yaml

dev-tools:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.10
	#go install github.com/favadi/protoc-go-inject-tag@latest

proto:
	cd proto && buf generate --template ./buf/buf.gen.local.yaml --config ./buf/buf.yaml .
	protoc-go-inject-tag -input="proto/gen/go/api/*.pb.go"

	protoc \
		--plugin=./web/node_modules/.bin/protoc-gen-ts_proto \
		--proto_path=./proto/src \
		--ts_proto_out=./web/src/lib/proto/ \
		--ts_proto_opt="esModuleInterop=true,forceLong=long" \
		./proto/src/**/*.proto	

tidy:
	go mod tidy -C .

build: $(foreach app,$(APPS),build-$(app))

build-%:
	go build -C ./cmd/$* -ldflags "-s -w" -o $(BUILD_OUTPUT)/$*$(EXE)

cp-yaml:
	$(foreach file,$(YAML_CFGS),cp app/$(file)/config.yaml $(BUILD_OUTPUT)/config.$(file).yaml;)
	cp app/common.yaml $(BUILD_OUTPUT)/config.common.yaml

clean:
	$(foreach app,$(APPS),rm $(BUILD_OUTPUT)/$(app)$(EXE);)
	$(foreach file,$(YAML_CFGS),rm $(BUILD_OUTPUT)/config.$(file).yaml;)

	#rm $(BUILD_OUTPUT)/config.common.yaml

.PHONY: proto dev-tools build cp-yaml clean tidy