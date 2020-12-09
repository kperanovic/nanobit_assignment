image: image-go-web image-go-worker

image-go-web:
	@echo "Building go web image..."
	@docker build -t go-web -f go/web/Dockerfile .

image-go-worker:
	@echo "Building go worker image..."
	@docker build -t go-worker -f go/worker/Dockerfile .
