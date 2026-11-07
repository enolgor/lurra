import path from 'node:path';
import tailwindcss from '@tailwindcss/vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import { defineConfig } from 'vite';

// https://vite.dev/config/
export default defineConfig({
	plugins: [tailwindcss(), svelte()],
	// $lib alias, also set in tsconfig.json and tsconfig.app.json
	resolve: {
		alias: { $lib: path.resolve('./src/lib') }
	},
	server: {
		// lets `npm run dev` talk to the go server without CORS
		proxy: { '/api': 'http://localhost:8080' }
	}
});
