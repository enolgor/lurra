import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { createRawSnippet } from 'svelte';
import { expect, test, vi } from 'vitest';
import MenuOption from './MenuOption.svelte';

function renderOption(current: boolean) {
	const onclick = vi.fn();
	render(MenuOption, {
		current,
		onclick,
		children: createRawSnippet(() => ({ render: () => '<span>Dark</span>' }))
	});
	return { onclick, button: screen.getByRole('button', { name: 'Dark' }) };
}

test('marks the current choice with a check and aria-current', () => {
	const { button } = renderOption(true);

	expect(button).toHaveAttribute('aria-current', 'true');
	expect(button.querySelector('.lucide-check')).toBeInTheDocument();
});

test('any other choice has neither', () => {
	const { button } = renderOption(false);

	expect(button).toHaveAttribute('aria-current', 'false');
	expect(button.querySelector('.lucide-check')).not.toBeInTheDocument();
});

test('a click chooses it', async () => {
	const { onclick, button } = renderOption(false);

	await userEvent.setup().click(button);

	expect(onclick).toHaveBeenCalledOnce();
});
