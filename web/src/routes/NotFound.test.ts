import { render, screen } from '@testing-library/svelte';
import { expect, test } from 'vitest';
import NotFound from './NotFound.svelte';

test('says the page does not exist and links back home', () => {
	render(NotFound);

	expect(screen.getByRole('heading', { level: 1, name: 'Page not found' })).toBeInTheDocument();
	expect(screen.getByText("There's nothing at this address.")).toBeInTheDocument();
	expect(screen.getByRole('link', { name: 'Back to your todos' })).toHaveAttribute('href', '/');
});
