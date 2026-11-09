import { render, screen } from '@testing-library/svelte';
import { expect, test } from 'vitest';
import Home from './Home.svelte';

test('has the page heading and the messages in the base locale (en)', () => {
	render(Home);

	expect(screen.getByRole('heading', { level: 1, name: 'Todos' })).toBeInTheDocument();
	expect(screen.getByText('Hello world ada')).toBeInTheDocument();
	expect(screen.getByText('My name is Ada Lovelace')).toBeInTheDocument();
});
