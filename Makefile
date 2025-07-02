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
	rm -f $(WIRE_DIR)/wire_gen.go

generate-mocks:
	@echo "generate mocks file
	mockgen -source=internal/domain/domain_order_core/repository.go -destination=internal/infrastructure/mocks/order_repository_mock.go -package=mocks
	mockgen -source=internal/domain/domain_product_core/service.go -destination=internal/infrastructure/mocks/product_service_mock.go -package=mocks 
	mockgen -source=internal/domain/domain_payment_core/repository.go -destination=internal/infrastructure/mocks/payment_repository_mock.go -package=mocks
	mockgen -source=internal/domain/domain_product_core/service.go -destination=internal/infrastructure/mocks/product_service_mock.go -package=mocks
	mockgen -source=internal/infrastructure/payment/payment_proxy.go -destination=internal/infrastructure/mocks/payment_proxy_mock.go -package=mocks