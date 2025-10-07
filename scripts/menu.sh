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
HIGHLIGHT="\033[1;30;48;2;166;227;161m"
CODEX_HIGHLIGHT_GREEN="\033[1;30;48;2;108;207;93m"
CODEX_HIGHLIGHT_PINK="\033[38;2;20;20;20;48;2;234;79;146m"

ENTER_KEY="[ ⏎ Enter ]"

# Read arrow keys
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

  # Gather Makefile targets
  entries=()

  # Debug: write to file
  {
    echo "=== DEBUG LOG ==="
    echo "Current directory: $(pwd)"
    echo "Makefile exists: $([ -f Makefile ] && echo 'yes' || echo 'no')"
    echo "makefile exists: $([ -f makefile ] && echo 'yes' || echo 'no')"
    echo ""
    echo "First 5 lines with ##:"
    grep -E ":.*##" [Mm]akefile 2>&1 | head -5
    echo ""
    echo "After first grep:"
    grep -E "^[a-zA-Z0-9_-]+:.*##" [Mm]akefile 2>&1 | head -5
    echo ""
    echo "After filtering menu:"
    grep -E "^[a-zA-Z0-9_-]+:.*##" [Mm]akefile | grep -vE "^_|^menu" | head -5
    echo ""
    echo "After sed:"
    grep -E "^[a-zA-Z0-9_-]+:.*##" [Mm]akefile | grep -vE "^_|^menu" | sed -E "s/^([a-zA-Z0-9_-]+):.*##[[:space:]]*(.*)/\1|\2/" | head -5
  } > /tmp/menu_debug.log 2>&1

  while IFS= read -r line; do
    target="${line%%|*}"
    desc="${line#*|}"
    entries+=("$target|$desc")
  done < <(
    grep -E "^[a-zA-Z0-9_-]+:.*##" [Mm]akefile 2>/dev/null |
      grep -vE "^_|^menu" |
      sed -E "s/^([a-zA-Z0-9_-]+):.*##[[:space:]]*(.*)/\1|\2/" |
      sort
  )

  # Debug: append results to file
  {
    echo ""
    echo "Found ${#entries[@]} entries:"
    for e in "${entries[@]}"; do
      echo "  - $e"
    done
  } >> /tmp/menu_debug.log 2>&1

  # Fallback for targets without descriptions
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

  # Separate options and descriptions
  options=("exit")
  descs=("quit this menu")
  for entry in "${entries[@]}"; do
    options+=("${entry%%|*}")
    descs+=("${entry#*|}")
  done

  selected=1
  max_index=$((${#options[@]} - 1))

  while true; do
    clear
    printf "${CODEX_PINK}---------------------------------------------${NC}\n"
    printf "${CODEX_GREEN}make commands for <codex>${NC}\n"
    printf "${CODEX_PINK}---------------------------------------------${NC}\n"

    for i in "${!options[@]}"; do
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
    '')
      break
      ;;
    q | Q)
      printf "\n${TEAL}Exiting menu...${NC}\n"
      exit 0
      ;;
    [0-9])
      if [ "$key" -ge 0 ] && [ "$key" -le "$max_index" ]; then
        selected=$key
        break
      else
        printf "\n${RED}⚠ invalid number (must be 0-%d)${NC}\n" "$max_index"
        sleep 1.5
      fi
      ;;
    *)
      ;;
    esac
  done

  # Run selected target
  if [ "$selected" -eq 0 ]; then
    printf "\n${TEAL}Exiting menu...${NC}\n"
    exit 0
  fi

  TARGET="${options[$selected]}"
  clear
  printf "${TEAL}Running target: ${CODEX_PINK}%s${NC}\n" "$TARGET"
  make "$TARGET"

  printf "\nPress ${CODEX_HIGHLIGHT_PINK}${ENTER_KEY}${NC} to return to menu..."
  read -r dummy
done
