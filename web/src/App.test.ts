import { render, screen } from '@testing-library/svelte';
import { expect, test } from 'vitest';
import App from './App.svelte';

test('renders the heading and button', () => {
	render(App);

	expect(screen.getByRole('heading', { name: 'lurra' })).toBeInTheDocument();
	expect(screen.getByRole('button', { name: 'Click me' })).toBeInTheDocument();
});
