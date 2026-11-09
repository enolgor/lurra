import { createRouter } from 'sv-router';
import Home from './routes/Home.svelte';

export const { p, navigate, isActive, preload, route } = createRouter({
	'/': Home,
	'*': Home // for now every other path too; a not-found page comes with the real app
});
