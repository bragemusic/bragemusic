package types

import (
	"errors"
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
)

const constantCSS = `
  /* not used, set to same as default */
  --field-background: var(--default);
  /* not used, guess its for placeholder text in forms */
  --field-placeholder: var(--default-foreground);
  /* color on focused objects. Set to same as accent */
  --focus: var(--accent);
  /* descriptive texts. Table headers, field description, list headers */
  --muted: var(--default-foreground);
  /* seems not to be used atm. Should probably be used as popup text colors */
  --overlay-foreground: var(--foreground);
  /* Seems not to be used. Set to same as border */
  --scrollbar: var(--border);
  /* Seems not to be used. Set to same as border */
  --segment: var(--border);
  /* Seems not to be used. Set to same as foreground */
  --segment-foreground: var(--foreground);
  /* Separators in tables, probably others. Set to same as border */
  --separator: var(--border);
  /* foregrouncd color on surface, set to same as foreground *NOT USED* */
  --surface-foreground: var(--foreground);
  /* not used, set to same as foreground */
  --surface-secondary-foreground: var(--foreground);
  /* not used, set to same as surface-secondary */
  --surface-tertiary: var(--surface-secondary);
  /* not used, set to same as foreground */
  --surface-tertiary-foreground: var(--foreground);
`

type ColorScheme string

const (
	ColorSchemeLight ColorScheme = "light"
	ColorSchemeDark  ColorScheme = "dark"
)

var validationRegex = regexp.MustCompile("^#(?:[0-9a-fA-F]{3}){1,2}$")

type HexColor string

func (h HexColor) validate(paramName, themeName string) error {
	valid := validationRegex.MatchString(string(h))
	if !valid {
		return fmt.Errorf("'%s' in color theme '%s' is not a valid hex color", paramName, themeName)
	}
	return nil
}

func (h HexColor) RGBA() color.RGBA {
	hex := strings.TrimPrefix(string(h), "#")

	if len(hex) != 6 {
		panic("expected 6 hex digits (RRGGBB)")
	}

	v, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		panic(err)
	}
	opacity := 1
	return color.RGBA{
		R: uint8(v >> 16),
		G: uint8(v >> 8),
		B: uint8(v),
		A: min(255, max(0, uint8(opacity*255))),
	}
}

func (h HexColor) NRGBA(opacity float32) color.NRGBA {
	hex := strings.TrimPrefix(string(h), "#")

	if len(hex) != 6 {
		panic("expected 6 hex digits (RRGGBB)")
	}

	v, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		panic(err)
	}

	return color.NRGBA{
		R: uint8(v >> 16),
		G: uint8(v >> 8),
		B: uint8(v),
		A: min(255, max(0, uint8(opacity*255))),
	}
}

type ThemeColors struct {
	Accent            HexColor `toml:"accent" desc:"Main brand color"`
	AccentForeground  HexColor `toml:"accent_foreground" desc:"Foreground color on components with accent background."`
	Background        HexColor `toml:"background" desc:"Background color of the main frame."`
	Border            HexColor `toml:"border" desc:"Border color."`
	Danger            HexColor `toml:"danger" desc:"Error message color."`
	DangerForeground  HexColor `toml:"danger_foreground" desc:"Foreground color on components with error background."`
	Default           HexColor `toml:"default" desc:"Used in many places. Its used for cancel buttons, hovers, background on fallback icons and background in input fields."`
	DefaultForeground HexColor `toml:"default_foreground" desc:"Used for player buttons, excluding play/pause."`
	FieldForeground   HexColor `toml:"field_foreground" desc:"Text color in forms."`
	Foreground        HexColor `toml:"foreground" desc:"Main text color."`
	Overlay           HexColor `toml:"overlay" desc:"Background in popup windows."`
	Success           HexColor `toml:"success" desc:"Success message color."`
	SuccessForeground HexColor `toml:"success_foreground" desc:"Foreground color on components with success background."`
	Surface           HexColor `toml:"surface" desc:"Background color of player, cards and notifications."`
	Surface2          HexColor `toml:"surface2" desc:"Background color of the menu."`
	Warning           HexColor `toml:"warning" desc:"Warning message color."`
	WarningForeground HexColor `toml:"warning_foreground" desc:"Foreground color on components with warning background."`
}

type Theme struct {
	Name        string      `toml:"name" desc:"The name of the theme. Used in the GUI for human readable selection."`
	ColorScheme ColorScheme `toml:"color_scheme" desc:"Define if your color scheme is light or dark. Nothing else is accepted."`
	Colors      ThemeColors `toml:"colors" desc:"Definition of all colors of the theme."`
}

type ThemeDescription struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (t Theme) Validate(themeName string) error {
	errs := []error{}

	if t.Name == "" {
		errs = append(errs, fmt.Errorf("'name' in theme '%s' is not set", themeName))
	}

	if t.ColorScheme != ColorSchemeLight && t.ColorScheme != ColorSchemeDark {
		errs = append(errs, fmt.Errorf("'color_scheme' in theme '%s' must be either light or dark", themeName))
	}

	errs = append(errs, t.Colors.Accent.validate("accent", themeName))
	errs = append(errs, t.Colors.AccentForeground.validate("accent_foreground", themeName))
	errs = append(errs, t.Colors.Background.validate("background", themeName))
	errs = append(errs, t.Colors.Border.validate("border", themeName))
	errs = append(errs, t.Colors.Danger.validate("danger", themeName))
	errs = append(errs, t.Colors.DangerForeground.validate("danger_foreground", themeName))
	errs = append(errs, t.Colors.Default.validate("default", themeName))
	errs = append(errs, t.Colors.DefaultForeground.validate("default_foreground", themeName))
	errs = append(errs, t.Colors.FieldForeground.validate("field_foreground", themeName))
	errs = append(errs, t.Colors.Foreground.validate("foreground", themeName))
	errs = append(errs, t.Colors.Overlay.validate("overlay", themeName))
	errs = append(errs, t.Colors.Success.validate("success", themeName))
	errs = append(errs, t.Colors.SuccessForeground.validate("success_foreground", themeName))
	errs = append(errs, t.Colors.Surface.validate("surface", themeName))
	errs = append(errs, t.Colors.Surface2.validate("surface2", themeName))
	errs = append(errs, t.Colors.Warning.validate("warning", themeName))
	errs = append(errs, t.Colors.WarningForeground.validate("warning_foreground", themeName))

	return errors.Join(errs...)
}

func (t Theme) CSS(themeName string) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, ".%s,\n", themeName)
	fmt.Fprintf(&sb, "[data-theme=\"%s\"] {\n", themeName)
	fmt.Fprintf(&sb, "  color-scheme: %s;\n", t.ColorScheme)

	fmt.Fprintf(&sb, "  --accent: %s;\n", t.Colors.Accent)
	fmt.Fprintf(&sb, "  --accent-foreground: %s;\n", t.Colors.AccentForeground)
	fmt.Fprintf(&sb, "  --background: %s;\n", t.Colors.Background)
	fmt.Fprintf(&sb, "  --border: %s;\n", t.Colors.Border)
	fmt.Fprintf(&sb, "  --danger: %s;\n", t.Colors.Danger)
	fmt.Fprintf(&sb, "  --danger-foreground: %s;\n", t.Colors.DangerForeground)
	fmt.Fprintf(&sb, "  --default: %s;\n", t.Colors.Default)
	fmt.Fprintf(&sb, "  --default-foreground: %s;\n", t.Colors.DefaultForeground)
	fmt.Fprintf(&sb, "  --field-foreground: %s;\n", t.Colors.FieldForeground)
	fmt.Fprintf(&sb, "  --foreground: %s;\n", t.Colors.Foreground)
	fmt.Fprintf(&sb, "  --overlay: %s;\n", t.Colors.Overlay)
	fmt.Fprintf(&sb, "  --success: %s;\n", t.Colors.Success)
	fmt.Fprintf(&sb, "  --success-foreground: %s;\n", t.Colors.SuccessForeground)
	fmt.Fprintf(&sb, "  --surface: %s;\n", t.Colors.Surface)
	fmt.Fprintf(&sb, "  --surface-secondary: %s;\n", t.Colors.Surface2)
	fmt.Fprintf(&sb, "  --warning: %s;\n", t.Colors.Warning)
	fmt.Fprintf(&sb, "  --warning-foreground: %s;\n", t.Colors.WarningForeground)

	fmt.Fprintf(&sb, "  %s\n", strings.TrimSpace(constantCSS))
	fmt.Fprint(&sb, "}\n")
	return sb.String()
}
