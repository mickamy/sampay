import { paraglideVitePlugin } from "@inlang/paraglide-js";
import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

const host = process.env.HOST;

export default defineConfig({
  plugins: [
    paraglideVitePlugin({
      project: "./project.inlang",
      outdir: "./app/paraglide",
    }),
    tailwindcss(),
    reactRouter(),
    tsconfigPaths(),
  ],
  server: {
    allowedHosts: host ? [host] : [],
  },
});
