TEAL=\033[38;2;148;226;213m
PEACH=\033[38;2;250;179;135m
GREEN=\033[38;2;166;227;161m
RED=\033[38;2;243;139;168m
CODEX_PINK=\033[38;2;234;79;146m
CODEX_GREEN=\033[38;2;108;207;93m
NC=\033[0m# No Color
HIGHLIGHT=\033[1;30;48;2;166;227;161m# bold black text on green background
CODEX_HIGHLIGHT_GREEN=\033[1;30;48;2;108;207;93m# bold black text on codex green background
CODEX_HIGHLIGHT_PINK	=\033[1;30;48;2;108;207;93m# bold black text on codex green background

NOWMS = go run tools/nowms.go

# Load saved values if .env exists
-include .env

menu:
	@./menu.sh

build: ## build the web server application
	@echo "$(CODEX_PINK)> building web server application...$(NC)"
	@START=$$($(NOWMS)); \
		go build -o bin/codex github.com/gary-norman/forum/cmd/server; \
		STOP=$$($(NOWMS)); \
		DIFF=$$((STOP - START)); \
		SEC=$$((DIFF / 1000)); \
		MS=$$((DIFF % 1000)); \
		echo "$(GREEN)✓ build complete!$(NC) in $${SEC}.$${MS}s"

run: ## run the web server application
	@echo "$(CODEX_PINK)> running web server application...$(NC)"
	@START=$$($(NOWMS)); \
		bin/codex; \
		STOP=$$($(NOWMS)); \
		DIFF=$$((STOP - START)); \
		HR=$$((DIFF / 3600000)); \
		MIN=$$(((DIFF / 60000) % 60)); \
		SEC=$$(((DIFF / 1000) % 60)); \
		MS=$$((DIFF % 1000)); \
		echo "$(GREEN)✓ server stopped$(NC) Uptime: $(CODEX_PINK)$${HR}h:$${MIN}m:$${SEC}.$${MS}s"
		printf "$(CODEX_PINK)---------------------------------------------$(NC)\n"; \

configure: ## configure Docker build and run options
	@bash -c '\
		printf "$(CODEX_PINK)---------------------------------------------$(NC)\n"; \
		printf "$(CODEX_GREEN)> configuring Docker build and run options...$(NC)\n"; \
		printf "$(CODEX_PINK)---------------------------------------------$(NC)\n"; \
		printf "Enter image name $(CODEX_GREEN)(leave blank for $(CODEX_PINK)samuishark/codex-v1.0)$(NC): "; \
		read IMAGE; \
		IMAGE=$${IMAGE:-samuishark/codex-v1.0}; \
		printf "Enter container name $(CODEX_GREEN)(leave blank for $(CODEX_PINK)codex)$(NC): "; \
		read CONTAINER; \
		CONTAINER=$${CONTAINER:-codex}; \
		printf "Enter local port number $(CODEX_GREEN)(leave blank for $(CODEX_PINK)8888)$(NC): "; \
		read PORT; \
		PORT=$${PORT:-8888}; \
		echo "IMAGE=$$IMAGE" > .env; \
		echo "CONTAINER=$$CONTAINER" >> .env; \
		echo "PORT=$$PORT" >> .env; \
		printf "$(GREEN)✓ configuration saved to $(CODEX_PINK).env$(NC)\n" \
		printf "$(CODEX_PINK)---------------------------------------------$(NC)\n"; \
	'

reset-config: ## reset Docker configuration by deleting .env
	@rm -f .env
	@printf "$(CODEX_GREEN)⚠ configuration reset$(NC) — run $(CODEX_PINK)'make configure'$(NC) to set new values\n"

build-image: ## build the Docker image
	@bash -c '\
		START=$$($(NOWMS)); \
		printf "$(CODEX_GREEN)> building Docker image $(NC)with tag: $(CODEX_PINK)%s$(NC)\n" "$(IMAGE)"; \
		docker image build -t $(IMAGE) .; \
		STOP=$$($(NOWMS)); \
		DIFF=$$((STOP - START)); \
		SEC=$$((DIFF / 1000)); \
		MS=$$((DIFF % 1000)); \
		printf "$(GREEN)✓ build complete!$(NC) in $(CODEX_PINK)%s.%03d$(NC)s\n" "$$SEC" "$$MS" \
	'

run-container: ## run the Docker container
	@bash -c '\
		START=$$($(NOWMS)); \
		printf "$(TEAL)> running Docker container $(NC)as $(CODEX_PINK)%s $(NC)from image: $(CODEX_PINK)%s $(NC)using port: $(CODEX_PINK)%s$(NC)\n" "$(CONTAINER)" "$(IMAGE)" "$(PORT)"; \
		docker run -d -p $(PORT):8888 --name $(CONTAINER) $(IMAGE); \
		STOP=$$($(NOWMS)); \
		DIFF=$$((STOP - START)); \
		SEC=$$((DIFF / 1000)); \
		MS=$$((DIFF % 1000)); \
		printf "$(GREEN)✓ task complete!$(NC)in $(CODEX_PINK)%s.%03d$(NC)s\n" "$$SEC" "$$MS" \
	'
