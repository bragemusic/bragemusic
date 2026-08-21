export const lightTheme = {
  // BaseColors
  background: '#F3F4F6',
  foreground: '#0B0D12',
  divider: '#E1E3E8',
  overlay: 'rgba(11,13,18,0.6)',
  focus: '#C89B3C',

  content1: '#0B0D12',
  content2: '#4B4F5E',
  content3: '#7A7F91',
  content4: '#B7BBC9',

  // Semantic colors
  default: '#1B1E27',

  primary: {
    50:  '#FBF6EA',
    100: '#F2E6C2',
    200: '#E8D69A',
    300: '#DEC671',
    400: '#D4B649',
    500: '#C89B3C', // DEFAULT
    600: '#A87F2F',
    700: '#876224',
    800: '#654718',
    900: '#432C0C',
    DEFAULT: '#C89B3C',
    foreground: '#0B0D12',
  },

  secondary: {
    50:  '#ECEFF6',
    100: '#D6DCEF',
    200: '#B0BBE0',
    300: '#8A9AD1',
    400: '#6479C2',
    500: '#3E5FB3', // DEFAULT
    600: '#344F94',
    700: '#2A3F75',
    800: '#202F56',
    900: '#161F37',
    DEFAULT: '#3E5FB3',
    foreground: '#F3F4F6',
  },

  success: {
    DEFAULT: '#2FA36B',
    foreground: '#07150F',
  },
  warning: {
    DEFAULT: '#D4A72C',
    foreground: '#1F1600',
  },
  danger: {
    DEFAULT: '#C24141',
    foreground: '#2B0A0A',
  },
};

export const darkTheme = {
  // BaseColors
  background: '#0B0D12',      // near-black with blue bias
  foreground: '#E6E8EE',      // soft white, not pure
  divider: '#1E2230',
  overlay: 'rgba(0,0,0,0.65)',
  focus: '#C89B3C',

  content1: '#E6E8EE',        // primary text
  content2: '#A3A8BC',        // secondary text
  content3: '#73788E',        // tertiary
  content4: '#4A4E5F',        // disabled

  highlight: '#E0B45A',
  popupbg: '#111420',

  // Semantic colors
  default: '#1B1E27',

  primary: {
    50:  '#2A2413',
    100: '#3A3118',
    200: '#534622',
    300: '#6C5B2D',
    400: '#857038',
    500: '#C89B3C', // DEFAULT (muted gold)
    600: '#D8B865',
    700: '#E4CD8F',
    800: '#EFE3BA',
    900: '#F7F1E1',
    DEFAULT: '#C89B3C',
    foreground: '#0B0D12',
  },

  secondary: {
    50:  '#0E1426',
    100: '#141C3A',
    200: '#1D2854',
    300: '#26356E',
    400: '#2F4288',
    500: '#3E5FB3', // DEFAULT (steel blue)
    600: '#5A7BC7',
    700: '#7896DA',
    800: '#A3BAEC',
    900: '#D4E1FA',
    DEFAULT: '#3E5FB3',
    foreground: '#0B0D12',
  },

  success: {
    DEFAULT: '#2FA36B',
    foreground: '#07150F',
  },
  warning: {
    DEFAULT: '#D4A72C',
    foreground: '#1F1600',
  },
  danger: {
    DEFAULT: '#C24141',
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
