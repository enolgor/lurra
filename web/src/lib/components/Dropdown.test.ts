import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { createRawSnippet } from 'svelte';
import { expect, test } from 'vitest';
import Dropdown from './Dropdown.svelte';
import { dropdown, trigger } from '$lib/test/dropdown';

// A menu with a title (not a choice) and one choice button.
function renderDropdown() {
	render(Dropdown, {
		trigger: createRawSnippet(() => ({ render: () => '<span>Open</span>' })),
		children: createRawSnippet(() => ({
			render: () => '<li><span>Title</span><button type="button">Choice</button></li>'
		}))
	});
	return userEvent.setup();
}

test('opens and closes on its button', async () => {
	const user = renderDropdown();

	await user.click(trigger());
	expect(dropdown()).toHaveAttribute('open');
	expect(screen.getByRole('button', { name: 'Choice' })).toBeVisible();

	await user.click(trigger());
	expect(dropdown()).not.toHaveAttribute('open');
});

test('closes when a choice is clicked, but not on other clicks inside', async () => {
	const user = renderDropdown();

	await user.click(trigger());
	await user.click(screen.getByText('Title'));
	expect(dropdown()).toHaveAttribute('open');

	await user.click(screen.getByRole('button', { name: 'Choice' }));
	expect(dropdown()).not.toHaveAttribute('open');
});

test('closes on a click outside', async () => {
	const user = renderDropdown();

	await user.click(trigger());
	await user.click(document.body);
	expect(dropdown()).not.toHaveAttribute('open');
});

test('closes on Escape and gives focus back to its button, ignoring other keys', async () => {
	const user = renderDropdown();

	// closed: Escape does nothing (focus stays where it is)
	await user.keyboard('{Escape}');
	expect(trigger()).not.toHaveFocus();

	await user.click(trigger());
	await user.keyboard('a');
	expect(dropdown()).toHaveAttribute('open');

	await user.keyboard('{Escape}');
	expect(dropdown()).not.toHaveAttribute('open');
	expect(trigger()).toHaveFocus();
});
