
export const lightTheme = {
  // BaseColors
  background: '#F4F4F8',
  foreground: '#0A0A12',
  divider: '#E6E7EB',
  overlay: 'rgba(10,10,18,0.6)',
  focus: '#0F766E',

  content1: '#0A0A12',
  content2: '#55576A',
  content3: '#8F92A3',
  content4: '#C9CAD2',

  // Semantic colors
  default: '#1A1B24',

  primary: {
    50:  '#E6F6F4',
    100: '#CDEDE9',
    200: '#9BDCD4',
    300: '#69CABE',
    400: '#3FB7A7',
    500: '#0F766E', // DEFAULT
    600: '#0C5F59',
    700: '#094843',
    800: '#06312D',
    900: '#031917',
    DEFAULT: '#0F766E',
    foreground: '#F4F4F8',
  },

  secondary: {
    50:  '#F1EEFB',
    100: '#DDD6F6',
    200: '#C4B5FD',
    300: '#A78BFA',
    400: '#7C3AED',
    500: '#6D28D9', // DEFAULT
    600: '#5B21B6',
    700: '#4C1D95',
    800: '#3B1574',
    900: '#2A0E53',
    DEFAULT: '#6D28D9',
    foreground: '#F4F4F8',
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

export const darkTheme = {
  // BaseColors
  background: '#0A0A12',
  foreground: '#F4F4F8',
  divider: '#25262D',
  overlay: 'rgba(0,0,0,0.6)',
  focus: '#0F766E',

  content1: '#F4F4F8',
  content2: '#A6A8B8',
  content3: '#8F92A3',
  content4: '#5B5C66',

  highlight: '#FF61C7',
  popupbg: '#10121B',

  // Semantic colors
  default: '#1A1B24',

  primary: {
    50:  '#0B2E2B',
    100: '#0D3B37',
    200: '#10504A',
    300: '#13645D',
    400: '#167871',
    500: '#0F766E', // DEFAULT
    600: '#1A8F86',
    700: '#22A89E',
    800: '#3FC4BA',
    900: '#7FE2DB',
    DEFAULT: '#0F766E',
    foreground: '#0A0A12',
  },

  secondary: {
    50:  '#1A1333',
    100: '#221A44',
    200: '#2E2260',
    300: '#3C2D80',
    400: '#4C37A3',
    500: '#6D28D9', // DEFAULT
    600: '#7C3AED',
    700: '#8B5CF6',
    800: '#A78BFA',
    900: '#C4B5FD',
    DEFAULT: '#6D28D9',
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

export const semanticBase = {
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
