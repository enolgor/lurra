import { createRouter } from 'sv-router';
import Home from '$routes/Home.svelte';
import NotFound from '$routes/NotFound.svelte';

// Only `p` (typed links) is used; createRouter also offers navigate, isActive, route, …
// eslint-disable-next-line @typescript-eslint/unbound-method -- p is a plain function, typed as a method
export const { p } = createRouter({
	'/': Home,
	'*': NotFound
});
