/// <reference types="vite-plugin-svgr/client" />
// import React from "react";
import ReactDOM from "react-dom/client";

import { WailsApi } from "./api/wailsApi.ts";
import { WailsPlayerApi } from "./api/wailsPlayerApi.ts";

import { App, ApiProvider } from "bragemusic";

function loadTheme() {
  const existing = document.getElementById("runtime-theme")
  if (existing) existing.remove()

  const link = document.createElement("link")
  link.id = "runtime-theme"
  link.rel = "stylesheet"
  link.href = "/config/custom.css"

  document.head.appendChild(link)
}

loadTheme()

ReactDOM.createRoot(document.getElementById("root")!).render(
  // <React.StrictMode>
  <ApiProvider api={new WailsApi()} playerApi={new WailsPlayerApi()}>
    <App />
  </ApiProvider>
  // </React.StrictMode>,
);
