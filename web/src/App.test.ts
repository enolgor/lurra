import { render, screen } from '@testing-library/svelte';
import { expect, test } from 'vitest';
import App from './App.svelte';

test.each(['/', '/todos/1'])('shows the current route for %s', async (path) => {
	history.replaceState(null, '', path);

	render(App);

	expect(await screen.findByText(`you are in route: ${path}`)).toBeInTheDocument();
	expect(screen.getByRole('heading', { name: 'lurra' })).toBeInTheDocument();
});
