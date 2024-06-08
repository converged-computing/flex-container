
# This assumes a build in the .devcontainer Dockerfile environment
FLUX_SCHED_ROOT ?= /opt/flux-sched
INSTALL_PREFIX ?= /usr
LIB_PREFIX ?= /usr/lib
LOCALBIN ?= $(shell pwd)/bin
COMMONENVVAR=GOOS=$(shell uname -s | tr A-Z a-z)
BUILDENVVAR=CGO_CFLAGS="-I${FLUX_SCHED_ROOT} -I${FLUX_SCHED_ROOT}/resource/reapi/bindings/c" CGO_LDFLAGS="-L${LIB_PREFIX} -L${LIB_PREFIX}/flux -L${FLUX_SCHED_ROOT}/resource/reapi/bindings -lreapi_cli -lflux-idset -lstdc++ -ljansson -lhwloc -lboost_system -lflux-hostlist -lboost_graph -lyaml-cpp"


.PHONY: all
all: flex

.PHONY: flex
flex: 
	go mod tidy
	mkdir -p ./bin
	$(COMMONENVVAR) $(BUILDENVVAR) go build -ldflags '-w' -o bin/flex-container src/cmd/main.go

.PHONY: clean
clean:
	rm -rf ./bin/server
