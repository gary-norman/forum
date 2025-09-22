// consoleThemes.js
import { palettes } from "./themes.js";

// Apply a theme to console.custom
export function applyCustomTheme(theme, flavor = null) {
  let palette;
  if (theme === "catppuccin") {
    if (!flavor)
      throw new Error(
        "Catppuccin requires a flavor: latte, frappe, macchiato, mocha",
      );
    palette = palettes.catppuccin[flavor];
  } else {
    palette = palettes[theme];
  }
  if (!palette) throw new Error("Unknown theme");

  const { error, expect, warn, info, bg } = palette;

  console.custom = {};

  function makeStyle(fg, bg) {
    return [
      `background-color:${bg}`,
      `color:${fg}`,
      "font-weight:bold",
      `border:1px solid ${fg}`,
      "padding:2px 6px",
      "border-radius:4px",
      "margin-right:4px",
      "white-space:nowrap",
    ].join("; ");
  }

  function wrap(origMethod, css, label) {
    return (...args) => {
      if (args.length === 0) return origMethod();
      const first = args[0];
      if (typeof first === "string") {
        if (first.includes("%c")) return origMethod(...args);
        return origMethod(`%c${first}`, css, ...args.slice(1));
      } else {
        return origMethod(`%c${label}`, css, ...args);
      }
    };
  }

  console.custom.log = wrap(
    console.log.bind(console),
    makeStyle(expect, bg),
    "[LOG]",
  );
  console.custom.warn = wrap(
    console.warn.bind(console),
    makeStyle(warn, bg),
    "[WARN]",
  );
  console.custom.error = wrap(
    console.error.bind(console),
    makeStyle(error, bg),
    "[ERROR]",
  );
  console.custom.info = wrap(
    console.info.bind(console),
    makeStyle(info, bg),
    "[INFO]",
  );
  console.custom.angryinfo = wrap(
    console.info.bind(console),
    makeStyle(info, bg),
    "[ANGRYINFO]",
  );

  // Mini palette preview
  printPalette(theme, flavor);
}

window.applyCustomTheme = applyCustomTheme;

// Print palette blocks for LOG/WARN/ERROR/INFO
export function printPalette(theme, flavor = null) {
  let palette;
  if (theme === "catppuccin") palette = palettes.catppuccin[flavor];
  else palette = palettes[theme];
  if (!palette) return;

  const { error, expect, warn, info, bg } = palette;
  const blocks = ["LOG", "WARN", "ERROR", "INFO", "ANGRYINFO"];
  const colors = [expect, warn, error, info, error];

  let formatStr = "";
  let styles = [];
  blocks.forEach((label, i) => {
    formatStr += `%c ${label} `;
    styles.push(
      `background-color:${bg}; color:${colors[i]}; font-weight:bold; border:1px solid ${colors[i]}; padding:2px 6px; border-radius:4px; margin-right:4px; white-space:nowrap; cursor:pointer;`,
    );
  });

  console.log(formatStr, ...styles);
}

// Flatten themes into a list for selection
export const themeList = [];
let idx = 1;
Object.entries(palettes).forEach(([themeName, val]) => {
  if (themeName === "catppuccin") {
    Object.keys(val).forEach((flavor) =>
      themeList.push({ index: idx++, theme: "catppuccin", flavor }),
    );
  } else {
    themeList.push({ index: idx++, theme: themeName });
  }
});

// Pick a theme by index
export function pickTheme(index) {
  const t = themeList.find((t) => t.index === index);
  if (!t) return console.warn("Invalid theme number");
  applyCustomTheme(t.theme, t.flavor);
  console.log(
    `%cTheme applied: ${t.theme}${t.flavor ? " - " + t.flavor : ""}`,
    "color:green; font-weight:bold;",
  );
}

// Attach to window for console use
window.pickTheme = pickTheme;

// Show all themes with clickable-looking blocks
export function showThemesClickable() {
  console.log("%cAvailable Themes", "font-weight:bold; font-size:14px;");
  themeList.forEach((t) => {
    const label =
      t.theme === "catppuccin" ? `Catppuccin - ${t.flavor}` : t.theme;
    const style =
      "background:#222; color:#fff; padding:2px 6px; border-radius:4px; cursor:pointer; font-weight:bold; margin-right:4px;";
    console.log(`%c[${t.index}] ${label}`, style);
    printPalette(t.theme, t.flavor);
  });
  console.log(
    "%cActivate a theme by calling pickTheme(number)",
    "font-style:italic; color:#aaa;",
  );
}

window.showThemesClickable = showThemesClickable;
