import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tsconfigPaths from "vite-tsconfig-paths";
import tailwindcss from "@tailwindcss/vite";
import svgr from "vite-plugin-svgr";
import path from "node:path";

const frontend = path.resolve(__dirname, "frontend");

export default defineConfig({
  root: frontend,

  server: {
    allowedHosts: ["picard.local"],
    proxy: {
      "/api": {
        target: "http://localhost:3000",
        changeOrigin: true,
        secure: false,
      },
      "/auth": {
        target: "http://localhost:3000",
        changeOrigin: true,
        secure: false,
      },
      "/config": {
        target: "http://localhost:3000",
        changeOrigin: true,
      },
      "/info": {
        target: "http://localhost:3000",
        changeOrigin: true,
      },
    },
  },

  plugins: [
    react(),
    tsconfigPaths(),
    tailwindcss(),
    svgr(),
  ],

  resolve: {
    alias: {
      "@": path.resolve(frontend, "lib"),
    },
  },
});
