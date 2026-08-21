/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./frontend/index.html",
    "./frontend/lib/layouts/**/*.{js,ts,jsx,tsx,mdx}",
    "./frontend/lib/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./frontend/lib/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./node_modules/@heroui/theme/dist/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      spacing: {
        body: "1400px",
      },
    },
  },
  darkMode: "class",
};
