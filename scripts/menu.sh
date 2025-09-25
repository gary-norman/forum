#!/usr/bin/env bash

# Color variables
TEAL="\033[38;2;148;226;213m"
PEACH="\033[38;2;250;179;135m"
GREEN="\033[38;2;166;227;161m"
RED="\033[38;2;243;139;168m"
YELLOW="\033[38;2;249;226;175m"
CODEX_PINK="\033[38;2;234;79;146m"
CODEX_GREEN="\033[38;2;108;207;93m"
NC="\033[0m"                                               # No Color
HIGHLIGHT="\033[1;30;48;2;166;227;161m"                    # bold black text on green background
CODEX_HIGHLIGHT_GREEN="\033[1;30;48;2;108;207;93m"         # bold black text on codex green background
CODEX_HIGHLIGHT_PINK="\033[38;2;20;20;20;48;2;234;79;146m" # bold black text on codex pink background

# Special characters
ENTER_KEY="[ ⏎ Enter ]"

# read arrow keys
read_arrow() {
  IFS= read -rsn1 key 2>/dev/null >&2
  if [[ $key == $'\x1b' ]]; then
    read -rsn2 key
    case $key in
    '[A') echo "up" ;;
    '[B') echo "down" ;;
    esac
  else
    echo "$key"
  fi
}

while true; do
  clear

  # gather Makefile targets
  entries=()
  while IFS= read -r line; do
    target="${line%%|*}"
    desc="${line#*|}"
    entries+=("$target|$desc")
  done < <(
    grep -E "^[a-zA-Z0-9_-]+:.*?##" Makefile |
      grep -vE "^_|^menu" |
      sed -E "s/^([a-zA-Z0-9_-]+):.*##[[:space:]]*(.*)/\1|\2/" |
      sort
  )

  # fallback for targets without descriptions
  if [ ${#entries[@]} -eq 0 ]; then
    while IFS= read -r line; do
      entries+=("$line|No description")
    done < <(
      grep -E "^[a-zA-Z0-9_-]+:([^=]|$)" Makefile |
        cut -d: -f1 |
        grep -vE "^_|^menu" |
        sort
    )
  fi

  # separate options and descriptions
  options=("exit")
  descs=("quit this menu")
  for entry in "${entries[@]}"; do
    options+=("${entry%%|*}")
    descs+=("${entry#*|}")
  done

  selected=1 # highlight first Makefile target initially
  max_index=$((${#options[@]} - 1))

  while true; do
    clear
    printf "${CODEX_PINK}---------------------------------------------${NC}\n"
    printf "${CODEX_GREEN}make commands for <codex>${NC}\n"
    printf "${CODEX_PINK}---------------------------------------------${NC}\n"

    for i in "${!options[@]}"; do
      # display numbers: 0 for exit, 1..n for targets
      if [ "$i" -eq 0 ]; then
        num=0
      else
        num=$i
      fi

      if [ "$i" -eq "$selected" ]; then
        printf "${CODEX_HIGHLIGHT_GREEN}%2d) %-15s - %s${NC}\n" "$num" "${options[$i]}" "${descs[$i]}"
      else
        printf "${CODEX_PINK}%2d)${CODEX_GREEN} %-15s${NC} - %s\n" "$num" "${options[$i]}" "${descs[$i]}"
      fi
    done

    printf "${CODEX_PINK}---------------------------------------------${NC}\n"
    printf "Use ↑/↓ to navigate, Enter to select, or type number (0-%d): " "$max_index"

    key=$(read_arrow)

    case "$key" in
    up)
      selected=$(((selected - 1 + ${#options[@]}) % ${#options[@]}))
      ;;
    down)
      selected=$(((selected + 1) % ${#options[@]}))
      ;;
    '') # Enter
      break
      ;;
    [0-9])
      if [ "$key" -ge 0 ] && [ "$key" -le "$max_index" ]; then
        selected=$key
        break
      else
        printf "${RED}⚠ invalid number${NC}\n"
        sleep 1
      fi
      ;;
    q)
      printf "${TEAL}Exiting menu...${NC}\n"
      exit 0
      ;;
    esac
  done

  # run selected target
  if [ "$selected" -eq 0 ]; then
    printf "${TEAL}Exiting menu...${NC}\n"
    exit 0
  fi

  TARGET="${options[$selected]}"
  clear
  printf "${TEAL}Running target: ${CODEX_PINK}%s${NC}\n" "$TARGET"
  make "$TARGET"
  printf "Press ${CODEX_HIGHLIGHT_PINK}${ENTER_KEY}${NC} to return to menu..."
  read dummy
done
