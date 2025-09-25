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
	@./scripts/menu.sh

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
	@./scripts/configure.sh

reset-config: ## reset Docker configuration by deleting .env and repopulating with dev/prod for db
	@./scripts/reset-config.sh

build-image: ## build the Docker image
	@bash -c '\
		START=$$($(NOWMS)); \
		printf "$(CODEX_GREEN)> building Docker image $(NC)with tag: $(CODEX_PINK)%s$(NC)\n" "$(IMAGE)"; \
		docker image build -t $(IMAGE) .; \
		STOP=$$($(NOWMS)); \
		DIFF=$$((STOP - START)); \
		SEC=$$((DIFF / 1000)); \
		MS=$$((DIFF % 1000)); \
		printf "$(GREEN)✓${CODEX_GREEN} build complete!$(NC) in $(CODEX_PINK)%s.%03d$(NC)s\n" "$$SEC" "$$MS" \
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
		printf "$(GREEN)✓${CODEX_GREEN} task complete!$(NC)in $(CODEX_PINK)%s.%03d$(NC)s\n" "$$SEC" "$$MS" \
	'
