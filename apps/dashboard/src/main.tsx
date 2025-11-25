import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./app";
import "@/shared/styles/globals.css";

// biome-ignore lint/style/noNonNullAssertion: false positive
createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>
);
