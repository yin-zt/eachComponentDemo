WINDOWS_APP_NAME=inner-redis.exe
APP_NAME=inner-redis
BUILD_DIR=build
BUILDDATE=`date "+%Y-%m-%d/%H:%M:%S"`
build = echo $(1) $(2)


mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))

CURRENT_DIR := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))
#ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

PROTO_SRC_PATH =${ROOT_DIR}/rpc


ALL_LINUX = linux-amd64 \
	linux-386 \
	linux-arm \
	linux-arm64

ALL = $(ALL_LINUX) \
		darwin-amd64 \
		darwin-arm64

build_all: $(ALL:%=build/%)

build/%:
	$(call build,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)))



#	$(call build,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)))
#	$(call md5,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)))
#	$(call tar,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)))
#	$(call delete,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)))

print:
	echo $(ALL:%=build/%)


clean:
	rmdir /S ${BUILD_DIR}

build_windows:
	go build -v -o $(BUILD_DIR)/$(WINDOWS_APP_NAME) ./main.go

vet:
	go vet main.go

lenv:
	go env -w CGO_ENABLED=0 GOOS=linux GOARCH=amd64

wenv:
	go env -w CGO_ENABLED=0 GOOS=windows GOARCH=amd64

build_linux:
	go build -v -o $(BUILD_DIR)/$(APP_NAME)  main.go

linux: lenv build_linux

windows: wenv build_windows

