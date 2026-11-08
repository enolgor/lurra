import { createRouter } from 'sv-router';
import Home from './routes/Home.svelte';

export const { p, navigate, isActive, preload, route } = createRouter({
	'*': Home
});
