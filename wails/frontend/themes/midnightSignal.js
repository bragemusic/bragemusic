export const lightTheme = {
  // BaseColors
  background: '#F5F7FA',
  foreground: '#0E1118',
  divider: '#E2E6ED',
  overlay: 'rgba(14,17,24,0.6)',
  focus: '#4FD1C5',

  content1: '#0E1118',
  content2: '#505566',
  content3: '#7C8194',
  content4: '#B8BCC9',

  // Semantic colors
  default: '#1E2230',

  primary: {
    50:  '#E6FBF9',
    100: '#C7F4EF',
    200: '#9EEAE2',
    300: '#75DFD6',
    400: '#5AD6CC',
    500: '#4FD1C5', // DEFAULT
    600: '#3FB2A8',
    700: '#2F8F88',
    800: '#206B67',
    900: '#104543',
    DEFAULT: '#4FD1C5',
    foreground: '#0E1118',
  },

  secondary: {
    50:  '#EEF2FF',
    100: '#D8E0FF',
    200: '#B4C0FF',
    300: '#8FA0FF',
    400: '#6A80FF',
    500: '#4C63E6', // DEFAULT
    600: '#3F52BF',
    700: '#324198',
    800: '#253071',
    900: '#181F4A',
    DEFAULT: '#4C63E6',
    foreground: '#F5F7FA',
  },

  success: {
    DEFAULT: '#22C59A',
    foreground: '#071510',
  },
  warning: {
    DEFAULT: '#E0A106',
    foreground: '#1F1600',
  },
  danger: {
    DEFAULT: '#D64545',
    foreground: '#2B0A0A',
  },
};

export const darkTheme = {
  // BaseColors
  background: '#0B0F16',     // deep blue-black
  foreground: '#E4E8F0',
  divider: '#1C2333',
  overlay: 'rgba(0,0,0,0.7)',
  focus: '#4FD1C5',

  content1: '#E4E8F0',
  content2: '#A0A6BC',
  content3: '#6F758A',
  content4: '#474B5C',

  highlight: '#67E8D5',
  popupbg: '#101622',

  // Semantic colors
  default: '#1E2230',

  primary: {
    50:  '#0E2A28',
    100: '#123835',
    200: '#184D49',
    300: '#1E635D',
    400: '#257973',
    500: '#4FD1C5', // DEFAULT (aqua-teal)
    600: '#6FE2D8',
    700: '#93EFE6',
    800: '#C1FAF4',
    900: '#E8FFFD',
    DEFAULT: '#4FD1C5',
    foreground: '#0B0F16',
  },

  secondary: {
    50:  '#121735',
    100: '#1A2150',
    200: '#26306F',
    300: '#32408E',
    400: '#3E50AD',
    500: '#4C63E6', // DEFAULT (electric indigo)
    600: '#6C7EF0',
    700: '#8E99F6',
    800: '#B8C1FB',
    900: '#E3E6FF',
    DEFAULT: '#4C63E6',
    foreground: '#0B0F16',
  },

  success: {
    DEFAULT: '#22C59A',
    foreground: '#071510',
  },
  warning: {
    DEFAULT: '#E0A106',
    foreground: '#1F1600',
  },
  danger: {
    DEFAULT: '#D64545',
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
