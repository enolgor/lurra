import { defineConfig } from "vite";

export default defineConfig({
  server: {
    // lets `npm run dev` talk to the go server without CORS
    proxy: { "/api": "http://localhost:8080" },
  },
});
