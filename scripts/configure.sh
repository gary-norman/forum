#!/usr/bin/env bash

# Color variables
TEAL="\033[38;2;148;226;213m"
PEACH="\033[38;2;250;179;135m"
GREEN="\033[38;2;166;227;161m"
RED="\033[38;2;243;139;168m"
YELLOW="\033[38;2;249;226;175m"
CODEX_PINK="\033[38;2;234;79;146m"
CODEX_GREEN="\033[38;2;108;207;93m"
NC="\033[0m"

printf "${CODEX_PINK}---------------------------------------------${NC}\n"
printf "${CODEX_GREEN}> configuring Database...${NC}\n"
printf "${CODEX_PINK}---------------------------------------------${NC}\n"
printf "${CODEX_GREEN}Select database environment:${NC}\n"
printf "${CODEX_PINK}1)${NC} development (SQLite)\n"
printf "${CODEX_PINK}2)${NC} production (SQLite)\n"
printf "Enter selection ${CODEX_PINK}[1-2]${NC}: "

read SELECTION

case $SELECTION in
1)
  DB_ENV="dev"
  DB_PATH="/var/lib/db-codex/dev_forum_database.db"
  ;;
2)
  DB_ENV="prod"
  DB_PATH="/var/lib/db-codex/forum_database.db"
  ;;
*)
  printf "${RED}✗ invalid selection, defaulting to development${NC}\n"
  DB_ENV="dev"
  DB_PATH="/var/lib/db-codex/dev_forum_database.db"
  ;;
esac

printf "${CODEX_PINK}---------------------------------------------${NC}\n"
printf "${CODEX_GREEN}> configuring Docker build and run options...${NC}\n"
printf "${CODEX_PINK}---------------------------------------------${NC}\n"
read -rp "Enter image name (default: samuishark/codex-v1.0): " IMAGE
IMAGE=${IMAGE:-samuishark/codex-v1.0}

read -rp "Enter container name (default: codex): " CONTAINER
CONTAINER=${CONTAINER:-codex}

read -rp "Enter local port number (default: 8888): " PORT
PORT=${PORT:-8888}

# Write everything fresh
cat >.env <<EOF
DB_ENV=$DB_ENV
DB_PATH=$DB_PATH
IMAGE=$IMAGE
CONTAINER=$CONTAINER
PORT=$PORT
EOF

printf "${GREEN}✓ configuration saved to ${CODEX_PINK}.env${NC}\n"
printf "${CODEX_PINK}---------------------------------------------${NC}\n"

printf "${YELLOW}⚠${CODEX_GREEN} configuration reset${NC} — run ${CODEX_PINK}'make configure'${NC} to set new values\n"
printf "${GREEN}✓${CODEX_GREEN} DB_ENV set to ${CODEX_PINK}%s${NC}\n" "$DB_PATH"
