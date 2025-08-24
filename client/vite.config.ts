import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

import { fileURLToPath } from "url";

// Fix for monaco workers
export default defineConfig({
  plugins: [react(), tailwindcss()],
  optimizeDeps: {
    include: ["@monaco-editor/react"],
  },
  resolve: {
    alias: {
      "monaco-editor": "monaco-editor/esm/vs/editor/editor.api",
      "@": path.resolve(path.dirname(fileURLToPath(import.meta.url)), "src"),
    },
  },
});
