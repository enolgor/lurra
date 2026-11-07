import { svelte } from '@sveltejs/vite-plugin-svelte';
import { defineConfig } from 'vite';

// https://vite.dev/config/
export default defineConfig({
	plugins: [svelte()],
	server: {
		// lets `npm run dev` talk to the go server without CORS
		proxy: { '/api': 'http://localhost:8080' }
	}
});
