.PHONY: all build run clean wire

BINARY_NAME=ddd_order_example

WIRE_DIR=internal/infrastructure/di

all: build

build: wire
	go build -o $(BINARY_NAME) main.go

wire:
	cd $(WIRE_DIR) && wire

run: build
	$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	#rm -f $(WIRE_DIR)/wire_gen.go