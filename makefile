
GREEN=\033[0;32m
BLUE=\033[1;34m
NC=\033[0m  # No Color

NOWMS = go run tools/nowms.go

build:
	@echo "$(BLUE)> building web server application...$(NC)"
	@START=$$($(NOWMS)); \
		go build -o bin/codex github.com/gary-norman/forum/cmd/server; \
		STOP=$$($(NOWMS)); \
		DIFF=$$((STOP - START)); \
		SEC=$$((DIFF / 1000)); \
		MS=$$((DIFF % 1000)); \
		echo "$(GREEN)✓ build complete!$(NC) in $${SEC}.$${MS}s"

run:
	@echo "$(BLUE)> running web server application...$(NC)"
	@START=$$($(NOWMS)); \
		bin/codex; \
		STOP=$$($(NOWMS)); \
		DIFF=$$((STOP - START)); \
		HR=$$((DIFF / 3600000)); \
		MIN=$$(((DIFF / 60000) % 60)); \
		SEC=$$(((DIFF / 1000) % 60)); \
		MS=$$((DIFF % 1000)); \
		echo "$(GREEN)✓ server stopped$(NC) Uptime: $${HR}h:$${MIN}m:$${SEC}.$${MS}s"
