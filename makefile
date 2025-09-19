
TEAL=\033[38;2;148;226;213m
PEACH=\033[38;2;250;179;135m
GREEN=\033[38;2;166;227;161m
NC=\033[0m  # No Color

NOWMS = go run tools/nowms.go

# Load saved values if .env exists
-include .env

build:
	@echo "$(TEAL)> building web server application...$(NC)"
	@START=$$($(NOWMS)); \
		go build -o bin/codex github.com/gary-norman/forum/cmd/server; \
		STOP=$$($(NOWMS)); \
		DIFF=$$((STOP - START)); \
		SEC=$$((DIFF / 1000)); \
		MS=$$((DIFF % 1000)); \
		echo "$(GREEN)✓ build complete!$(NC) in $${SEC}.$${MS}s"

run:
	@echo "$(TEAL)> running web server application...$(NC)"
	@START=$$($(NOWMS)); \
		bin/codex; \
		STOP=$$($(NOWMS)); \
		DIFF=$$((STOP - START)); \
		HR=$$((DIFF / 3600000)); \
		MIN=$$(((DIFF / 60000) % 60)); \
		SEC=$$(((DIFF / 1000) % 60)); \
		MS=$$((DIFF % 1000)); \
		echo "$(GREEN)✓ server stopped$(NC) Uptime: $(PEACH)$${HR}h:$${MIN}m:$${SEC}.$${MS}s"

configure:
	@bash -c '\
		read -p "Enter image name (leave blank for samuishark/codex-v1.0): " IMAGE; \
		IMAGE=$${IMAGE:-samuishark/codex-v1.0}; \
		read -p "Enter container name (leave blank for codex): " CONTAINER; \
		CONTAINER=$${CONTAINER:-codex}; \
		read -p "Enter local port number (leave blank for 8888): " PORT; \
		PORT=$${PORT:-8888}; \
		echo "IMAGE=$$IMAGE" > .env; \
		echo "CONTAINER=$$CONTAINER" >> .env; \
		echo "PORT=$$PORT" >> .env; \
		echo "$(GREEN)✓ configuration saved to .env$(NC)" \
	'

reset-config:
	@rm -f .env
	@printf "$(PEACH)⚠ configuration reset — run 'make configure' to set new values$(NC)\n"

build-image:
	@bash -c '\
		START=$$($(NOWMS)); \
		printf "$(TEAL)> building Docker image $(NC)with tag: $(PEACH)%s$(NC)\n" "$(IMAGE)"; \
		docker image build -t $(IMAGE) .; \
		STOP=$$($(NOWMS)); \
		DIFF=$$((STOP - START)); \
		SEC=$$((DIFF / 1000)); \
		MS=$$((DIFF % 1000)); \
		printf "$(GREEN)✓ build complete!$(NC) in $(PEACH)%s.%03d$(NC)s\n" "$$SEC" "$$MS" \
	'

run-container:
	@bash -c '\
		START=$$($(NOWMS)); \
		printf "$(TEAL)> running Docker container $(NC)as $(PEACH)%s $(NC)from image: $(PEACH)%s $(NC)using port: $(PEACH)%s$(NC)\n" "$(CONTAINER)" "$(IMAGE)" "$(PORT)"; \
		docker run -d -p $(PORT):8888 --name $(CONTAINER) $(IMAGE); \
		STOP=$$($(NOWMS)); \
		DIFF=$$((STOP - START)); \
		SEC=$$((DIFF / 1000)); \
		MS=$$((DIFF % 1000)); \
		printf "$(GREEN)✓ task complete!$(NC)in $(PEACH)%s.%03d$(NC)s\n" "$$SEC" "$$MS" \
	'
