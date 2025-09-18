// Package colors provides color palettes and ANSI escape codes for terminal output.
package colors

import "fmt"

// Convert hex color to ANSI 24-bit escape
func ansi(hex string) string {
	var r, g, b uint8
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

type ColorSet struct {
	// Accent colors
	Rosewater, Flamingo, Pink, Mauve, Red, Maroon, Peach, Yellow, Green, Teal, Sky, Sapphire, Blue, Lavender string

	// Neutrals
	Text, Subtext1, Subtext0, Overlay2, Overlay1, Overlay0 string
	Surface2, Surface1, Surface0, Base, Mantle, Crust      string

	// Utility
	Reset string
}

var palettes = map[string]*ColorSet{
	"Latte": {
		Rosewater: ansi("#dc8a78"), Flamingo: ansi("#dd7878"), Pink: ansi("#ea76cb"), Mauve: ansi("#8839ef"),
		Red: ansi("#d20f39"), Maroon: ansi("#e64553"), Peach: ansi("#fe640b"), Yellow: ansi("#df8e1d"),
		Green: ansi("#40a02b"), Teal: ansi("#179299"), Sky: ansi("#04a5e5"), Sapphire: ansi("#209fb5"),
		Blue: ansi("#1e66f5"), Lavender: ansi("#7287fd"),

		Text: ansi("#4c4f69"), Subtext1: ansi("#5c5f77"), Subtext0: ansi("#6c6f85"),
		Overlay2: ansi("#7c7f93"), Overlay1: ansi("#8c8fa1"), Overlay0: ansi("#9ca0b0"),
		Surface2: ansi("#acb0be"), Surface1: ansi("#bcc0cc"), Surface0: ansi("#ccd0da"),
		Base: ansi("#eff1f5"), Mantle: ansi("#e6e9ef"), Crust: ansi("#dce0e8"),
		Reset: "\033[0m",
	},

	"Frappe": {
		Rosewater: ansi("#f2d5cf"), Flamingo: ansi("#eebebe"), Pink: ansi("#f4b8e4"), Mauve: ansi("#ca9ee6"),
		Red: ansi("#e78284"), Maroon: ansi("#ea999c"), Peach: ansi("#ef9f76"), Yellow: ansi("#e5c890"),
		Green: ansi("#a6d189"), Teal: ansi("#81c8be"), Sky: ansi("#99d1db"), Sapphire: ansi("#85c1dc"),
		Blue: ansi("#8caaee"), Lavender: ansi("#babbf1"),

		Text: ansi("#c6d0f5"), Subtext1: ansi("#b5bfe2"), Subtext0: ansi("#a5adce"),
		Overlay2: ansi("#949cbb"), Overlay1: ansi("#838ba7"), Overlay0: ansi("#737994"),
		Surface2: ansi("#626880"), Surface1: ansi("#51576d"), Surface0: ansi("#414559"),
		Base: ansi("#303446"), Mantle: ansi("#292c3c"), Crust: ansi("#232634"),
		Reset: "\033[0m",
	},

	"Macchiato": {
		Rosewater: ansi("#f4dbd6"), Flamingo: ansi("#f0c6c6"), Pink: ansi("#f5bde6"), Mauve: ansi("#c6a0f6"),
		Red: ansi("#ed8796"), Maroon: ansi("#ee99a0"), Peach: ansi("#f5a97f"), Yellow: ansi("#eed49f"),
		Green: ansi("#a6da95"), Teal: ansi("#8bd5ca"), Sky: ansi("#91d7e3"), Sapphire: ansi("#7dc4e4"),
		Blue: ansi("#8aadf4"), Lavender: ansi("#b7bdf8"),

		Text: ansi("#cad3f5"), Subtext1: ansi("#b8c0e0"), Subtext0: ansi("#a5adcb"),
		Overlay2: ansi("#939ab7"), Overlay1: ansi("#8087a2"), Overlay0: ansi("#6e738d"),
		Surface2: ansi("#5b6078"), Surface1: ansi("#494d64"), Surface0: ansi("#363a4f"),
		Base: ansi("#24273a"), Mantle: ansi("#1e2030"), Crust: ansi("#181926"),
		Reset: "\033[0m",
	},

	"Mocha": {
		Rosewater: ansi("#f5e0dc"), Flamingo: ansi("#f2cdcd"), Pink: ansi("#f5c2e7"), Mauve: ansi("#cba6f7"),
		Red: ansi("#f38ba8"), Maroon: ansi("#eba0ac"), Peach: ansi("#fab387"), Yellow: ansi("#f9e2af"),
		Green: ansi("#a6e3a1"), Teal: ansi("#94e2d5"), Sky: ansi("#89dceb"), Sapphire: ansi("#74c7ec"),
		Blue: ansi("#89b4fa"), Lavender: ansi("#b4befe"),

		Text: ansi("#cdd6f4"), Subtext1: ansi("#bac2de"), Subtext0: ansi("#a6adc8"),
		Overlay2: ansi("#9399b2"), Overlay1: ansi("#7f849c"), Overlay0: ansi("#6c7086"),
		Surface2: ansi("#585b70"), Surface1: ansi("#45475a"), Surface0: ansi("#313244"),
		Base: ansi("#1e1e2e"), Mantle: ansi("#181825"), Crust: ansi("#11111b"),
		Reset: "\033[0m",
	},
}

var Colors *ColorSet

// UseFlavor sets the global Colors variable to a flavor
func UseFlavor(flavor string) *ColorSet {
	if c, ok := palettes[flavor]; ok {
		Colors = c
	} else {
		Colors = palettes["Mocha"]
	}
	return Colors
}
