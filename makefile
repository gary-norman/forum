
TEAL=\033[38;2;148;226;213m
PEACH=\033[38;2;250;179;135m
GREEN=\033[38;2;166;227;161m
RED=\033[38;2;243;139;168m
NC=\033[0m  # No Color

1="build"
2="run"
3="configure"
4="reset-config"
5="build-image"
6="run-container"

NOWMS = go run tools/nowms.go

# Load saved values if .env exists
-include .env

menu:
	@bash -c '\
		trap "echo; printf \"$(PEACH)Exiting menu...$(NC)\n\"; exit 0" INT; \
		while true; do \
			clear; \
			printf "$(TEAL)make commands for <codex>\n"; \
			printf "$(PEACH)---------------------------------------------$(NC)\n"; \
			if [ -f .env ]; then \
				printf "$(TEAL)> loaded configuration from .env$(NC)\n"; \
			else \
				printf "$(PEACH)⚠ no .env file found — run make configure$(NC)\n"; \
			fi; \
			printf "$(TEAL)---------------------------------------------$(NC)\n"; \
			options=("build\n" "run" "configure" "reset-config" "build-image" "run-container" "exit"); \
			PS3="$(PEACH)Enter command number (1-$${#options[@]}): $(NC)"; \
			select opt in "$${options[@]}"; do \
				if [ "$$opt" = "exit" ]; then \
					printf "$(PEACH)Exiting menu...$(NC)\n"; \
					exit 0; \
				elif [ -n "$$opt" ]; then \
					$(MAKE) $$opt; \
					printf "$(TEAL)---------------------------------------------$(NC)\n"; \
					read -p "$(PEACH)Press Enter to return to menu...$(NC)" dummy; \
					break; \
				else \
					printf "$(RED)⚠ invalid choice — please enter a number between 1 and $${#options[@]}$(NC)\n"; \
				fi; \
			done; \
		done'

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
		printf "$(TEAL)> configuring Docker build and run options...$(NC)\n"; \
		printf "$(TEAL)---------------------------------------------$(NC)\n"; \
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
