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

# Get current time in milliseconds
get_time_ms() {
  if command -v python3 &> /dev/null; then
    python3 -c 'import time; print(int(time.time() * 1000))'
  elif command -v node &> /dev/null; then
    node -e 'console.log(Date.now())'
  else
    # Fallback to seconds
    date +%s
  fi
}

# Format duration from milliseconds
format_duration() {
  local duration_ms=$1
  local hours=$((duration_ms / 3600000))
  local minutes=$(( (duration_ms % 3600000) / 60000 ))
  local seconds=$(( (duration_ms % 60000) / 1000 ))
  local ms=$((duration_ms % 1000))

  if [ $hours -gt 0 ]; then
    printf "%dh %dm %d.%03ds" "$hours" "$minutes" "$seconds" "$ms"
  elif [ $minutes -gt 0 ]; then
    printf "%dm %d.%03ds" "$minutes" "$seconds" "$ms"
  else
    printf "%d.%03ds" "$seconds" "$ms"
  fi
}

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

# Show menu function
show_menu() {
  local title="$1"
  local options_var="$2"
  local descs_var="$3"

  echo "DEBUG: Inside show_menu" >&2
  echo "DEBUG: options_var=$options_var, descs_var=$descs_var" >&2

  # Use eval to access array indirectly (compatible with bash 3.2)
  eval "local options=(\"\${${options_var}[@]}\")"
  eval "local descs=(\"\${${descs_var}[@]}\")"

  echo "DEBUG: After eval, options count: ${#options[@]}" >&2
  echo "DEBUG: Options: ${options[@]}" >&2

  local selected=1
  local max_index=$((${#options[@]} - 1))

  echo "DEBUG: Entering menu loop" >&2

  while true; do
    echo "DEBUG: Top of loop, about to clear" >&2
    clear
    printf "${CODEX_PINK}---------------------------------------------${NC}\n"
    printf "${CODEX_GREEN}%s${NC}\n" "$title"
    printf "${CODEX_PINK}---------------------------------------------${NC}\n"

    for i in "${!options[@]}"; do
      if [ "$i" -eq 0 ]; then
        num=0
      else
        num=$i
      fi

      if [ "$i" -eq "$selected" ]; then
        printf "${CODEX_HIGHLIGHT_GREEN}%2d) %-20s - %s${NC}\n" "$num" "${options[$i]}" "${descs[$i]}"
      else
        printf "${CODEX_PINK}%2d)${CODEX_GREEN} %-20s${NC} - %s\n" "$num" "${options[$i]}" "${descs[$i]}"
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
      echo "${options[$selected]}"
      return 0
      ;;
    q | Q)
      echo "exit"
      return 0
      ;;
    [0-9])
      if [ "$key" -ge 0 ] && [ "$key" -le "$max_index" ]; then
        echo "${options[$key]}"
        return 0
      else
        printf "\n${RED}⚠ invalid number (must be 0-%d)${NC}\n" "$max_index"
        sleep 1.5
      fi
      ;;
    *)
      ;;
    esac
  done
}

# Docker submenu
docker_menu() {
  local docker_options=("back" "configure" "reset-config" "build-image" "run-container")
  local docker_descs=("return to main menu" "configure Docker options" "reset configuration" "build Docker image" "run Docker container")

  while true; do
    choice=$(show_menu "Docker Commands" docker_options docker_descs)

    case "$choice" in
      back|exit)
        return
        ;;
      *)
        clear
        printf "${TEAL}Running target: ${CODEX_PINK}%s${NC}\n\n" "$choice"

        START=$(get_time_ms)
        make "$choice"
        EXIT_CODE=$?
        END=$(get_time_ms)
        DURATION=$((END - START))

        printf "\n${CODEX_PINK}---------------------------------------------${NC}\n"
        if [ $EXIT_CODE -eq 0 ]; then
          printf "${GREEN}✓ Task completed${NC} in ${CODEX_PINK}%s${NC}\n" "$(format_duration $DURATION)"
        else
          printf "${RED}✗ Task failed${NC} (exit code: %d) after ${CODEX_PINK}%s${NC}\n" "$EXIT_CODE" "$(format_duration $DURATION)"
        fi
        printf "${CODEX_PINK}---------------------------------------------${NC}\n"

        printf "\nPress ${CODEX_HIGHLIGHT_PINK}${ENTER_KEY}${NC} to return to menu..."
        read -r dummy
        ;;
    esac
  done
}

# Scripts submenu
scripts_menu() {
  local scripts_options=("back" "install-scripts" "verify-scripts" "backup-scripts")
  local scripts_descs=("return to main menu" "install/update scripts" "verify checksums" "backup scripts")

  while true; do
    choice=$(show_menu "Script Management" scripts_options scripts_descs)

    case "$choice" in
      back|exit)
        return
        ;;
      *)
        clear
        printf "${TEAL}Running target: ${CODEX_PINK}%s${NC}\n\n" "$choice"

        START=$(get_time_ms)
        make "$choice"
        EXIT_CODE=$?
        END=$(get_time_ms)
        DURATION=$((END - START))

        printf "\n${CODEX_PINK}---------------------------------------------${NC}\n"
        if [ $EXIT_CODE -eq 0 ]; then
          printf "${GREEN}✓ Task completed${NC} in ${CODEX_PINK}%s${NC}\n" "$(format_duration $DURATION)"
        else
          printf "${RED}✗ Task failed${NC} (exit code: %d) after ${CODEX_PINK}%s${NC}\n" "$EXIT_CODE" "$(format_duration $DURATION)"
        fi
        printf "${CODEX_PINK}---------------------------------------------${NC}\n"

        printf "\nPress ${CODEX_HIGHLIGHT_PINK}${ENTER_KEY}${NC} to return to menu..."
        read -r dummy
        ;;
    esac
  done
}

# Main menu loop
while true; do
  # Build main menu
  main_options=("exit" "build" "run" "🐳 Docker" "📜 Scripts")
  main_descs=("quit this menu" "build the application" "run the application" "Docker management" "script management")

  # Debug
  echo "DEBUG: About to call show_menu" >&2
  echo "DEBUG: Options: ${main_options[@]}" >&2

  choice=$(show_menu "make commands for <codex>" main_options main_descs)

  echo "DEBUG: Choice returned: $choice" >&2

  case "$choice" in
    exit)
      printf "\n${TEAL}Exiting menu...${NC}\n"
      exit 0
      ;;
    "🐳 Docker")
      docker_menu
      ;;
    "📜 Scripts")
      scripts_menu
      ;;
    *)
      clear
      printf "${TEAL}Running target: ${CODEX_PINK}%s${NC}\n\n" "$choice"

      START=$(get_time_ms)
      make "$choice"
      EXIT_CODE=$?
      END=$(get_time_ms)
      DURATION=$((END - START))

      printf "\n${CODEX_PINK}---------------------------------------------${NC}\n"
      if [ $EXIT_CODE -eq 0 ]; then
        printf "${GREEN}✓ Task completed${NC} in ${CODEX_PINK}%s${NC}\n" "$(format_duration $DURATION)"
      else
        printf "${RED}✗ Task failed${NC} (exit code: %d) after ${CODEX_PINK}%s${NC}\n" "$EXIT_CODE" "$(format_duration $DURATION)"
      fi
      printf "${CODEX_PINK}---------------------------------------------${NC}\n"

      printf "\nPress ${CODEX_HIGHLIGHT_PINK}${ENTER_KEY}${NC} to return to menu..."
      read -r dummy
      ;;
  esac
done
