COMPILER_BINARY_NAME=datagenc
TRANSPILED_BINARY_NAME=datagen

.PHONY: all build-compiler run clean

all: target build-compiler transpile run

build-compiler:
	go build -o target/$(COMPILER_BINARY_NAME) .

clean:
	go clean
	rm -f target/$(COMPILER_BINARY_NAME)
	rm -f target/$(TRANSPILED_BINARY_NAME)

target:
	mkdir target

transpile: build-compiler
	@echo "Running $(COMPILER_BINARY_NAME) with in=$(in) out=$(out)"
	target/$(COMPILER_BINARY_NAME) --in=$(in) --out=$(out)
	go build -C $(out) -o $(TRANSPILED_BINARY_NAME)
	cp $(out)/$(TRANSPILED_BINARY_NAME) target/$(TRANSPILED_BINARY_NAME)

run:
	./target/$(TRANSPILED_BINARY_NAME) -gen $(model):$(count)
%:
	@:
