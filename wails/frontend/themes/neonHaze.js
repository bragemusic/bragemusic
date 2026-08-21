// import { SemanticBaseColors, ThemeColors } from "@heroui/theme";

// neonHazeTheme.ts
export const lightTheme = {
  // BaseColors
  background: '#F4F4F8', // page background
  foreground: '#0A0A12', // main text color (in light mode)
  divider: '#E6E7EB',
  overlay: 'rgba(10,10,18,0.6)', // modal overlays (string ok)
  focus: '#007AFF', // focus ring uses primary blue
  content1: '#0A0A12', // primary text
  content2: '#55576A', // secondary text (between surface and primary)
  content3: '#8F92A3', // tertiary (meta)
  content4: '#C9CAD2', // disabled / subtle

  // Semantic colors
  default: '#1A1B24', // neutral surface used as "default" element color
  primary: {
    50: '#E6F0FF',
    100: '#BFE0FF',
    200: '#99D0FF',
    300: '#66B8FF',
    400: '#339FFF',
    500: '#007AFF', // DEFAULT
    600: '#0066E6',
    700: '#004FB3',
    800: '#00377A',
    900: '#001F3F',
    DEFAULT: '#007AFF',
    foreground: '#F4F4F8',
  },
  secondary: {
    50: '#F0EEFF',
    100: '#DED8FF',
    200: '#CDBFFF',
    300: '#B29AFF',
    400: '#9B7CFF',
    500: '#8B5CF6', // DEFAULT
    600: '#7A46E6',
    700: '#5B2EB3',
    800: '#3C1E7A',
    900: '#1E0F3F',
    DEFAULT: '#8B5CF6',
    foreground: '#F4F4F8',
  },
  success: {
    DEFAULT: '#22C55E',
    foreground: '#07150B',
  },
  warning: {
    DEFAULT: '#F59E0B',
    foreground: '#201300',
  },
  danger: {
    DEFAULT: '#EF4444',
    foreground: '#2B0A0A',
  },
};

export const darkTheme = {
  // BaseColors
  background: '#0A0A12', // deep dark background
  foreground: '#F4F4F8', // main text color (in dark mode)
  divider: '#25262D',
  overlay: 'rgba(0,0,0,0.6)',
  focus: '#8B5CF6', // purple focus on dark feels nicer
  content1: '#F4F4F8', // primary text
  content2: '#A6A8B8', // between surface and primary text (hex)
  content3: '#8F92A3',
  content4: '#5B5C66', // disabled / subtle
  highlight: '#FF61C7',
  popupbg: '#141622',

  // Semantic colors
  default: '#1A1B24',
  primary: {
    50: '#0E2A64',
    100: '#0C265B',
    200: '#0A204C',
    300: '#071836',
    400: '#005BFF',
    500: '#007AFF', // DEFAULT
    600: '#2A6DF0',
    700: '#1854C2',
    800: '#123C8A',
    900: '#091E48',
    DEFAULT: '#007AFF',
    foreground: '#0A0A12',
  },
  secondary: {
    50: '#140F2A',
    100: '#211737',
    200: '#3B2A5F',
    300: '#563F90',
    400: '#7C5CF0',
    500: '#8B5CF6', // DEFAULT
    600: '#7A46E6',
    700: '#5B2EB3',
    800: '#3C1E7A',
    900: '#1E0F3F',
    DEFAULT: '#8B5CF6',
    foreground: '#0A0A12',
  },
  success: {
    DEFAULT: '#16A34A',
    foreground: '#07150B',
  },
  warning: {
    DEFAULT: '#D97706',
    foreground: '#201300',
  },
  danger: {
    DEFAULT: '#DC2626',
    foreground: '#2B0A0A',
  },
};

// SemanticBaseColors (just the BaseColors portion for convenience)
export const semanticBase= {
  light: {
    background: lightTheme.background,
    foreground: lightTheme.foreground,
    divider: lightTheme.divider,
    overlay: lightTheme.overlay,
    focus: lightTheme.focus,
    content1: lightTheme.content1,
    content2: lightTheme.content2,
    content3: lightTheme.content3,
    content4: lightTheme.content4,
  },
  dark: {
    background: darkTheme.background,
    foreground: darkTheme.foreground,
    divider: darkTheme.divider,
    overlay: darkTheme.overlay,
    focus: darkTheme.focus,
    content1: darkTheme.content1,
    content2: darkTheme.content2,
    content3: darkTheme.content3,
    content4: darkTheme.content4,
    popupbg: darkTheme.popupbg,
  },
};
