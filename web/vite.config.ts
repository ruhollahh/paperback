/// <reference types="vite" />
import { defineConfig } from "vite";

export default defineConfig({
  build: {
    rollupOptions: {
      input: {
        main: "src/main.ts",
        "components/cart": "src/components/cart.ts",
        "components/other": "src/components/other.ts",
      },
      output: {
        entryFileNames: "[name].js",
        dir: "./static/dist",
      },
    },
  },
});
