/// <reference types="vite-plugin-svgr/client" />
// import React from "react";
import ReactDOM from "react-dom/client";

import App from "../lib/App.tsx";
import "@/styles/globals.css";
import { ApiProvider } from "../lib/api/ApiContext.tsx";
import { ServerApi } from "../lib/api/serverApi.ts";

function loadTheme() {
  const existing = document.getElementById("runtime-theme");

  if (existing) existing.remove();

  const link = document.createElement("link");

  link.id = "runtime-theme";
  link.rel = "stylesheet";
  link.href = "/config/custom.css";

  document.head.appendChild(link);
}

loadTheme();

const api = new ServerApi()

ReactDOM.createRoot(document.getElementById("root")!).render(
  // <React.StrictMode>
  <ApiProvider api={api} playerApi={api}>
    <App />
  </ApiProvider>,
  // </React.StrictMode>,
);
