import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tsconfigPaths from "vite-tsconfig-paths";
import tailwindcss from "@tailwindcss/vite";
import svgr from 'vite-plugin-svgr';

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:31145',
        changeOrigin: true,
      },
      '/config': {
        target: 'http://localhost:31145',
        changeOrigin: true,
      }
    }
  },
  plugins: [react(), tsconfigPaths(), tailwindcss(), svgr()],
});
