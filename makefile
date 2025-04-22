# ANSI color codes
GREEN=\033[0;32m
BLUE=\033[1;34m
NC=\033[0m  # No Color

build:
	@echo "$(BLUE)> building web server application...$(NC)"
	@go build -o bin/codex github.com/gary-norman/forum/cmd/server
	@echo "$(GREEN)✓ build complete!$(NC)"

run:
	@echo "$(BLUE)> running web server application...$(NC)"
	@bin/codex
	@echo "$(GREEN)✓ server stopped$(NC)"
