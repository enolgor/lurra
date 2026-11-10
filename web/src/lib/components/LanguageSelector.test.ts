import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { expect, test } from 'vitest';
import LanguageSelector from './LanguageSelector.svelte';
import { dropdown, trigger } from '$lib/test/dropdown';

const option = (name: string) => screen.getByRole('button', { name: new RegExp(name) });

test('lists every locale in its own language, with the current one marked', async () => {
	const user = userEvent.setup();
	render(LanguageSelector);

	await user.click(trigger());

	expect(dropdown()).toHaveAttribute('open');
	expect(screen.getByText('Select language')).toBeVisible();
	expect(
		screen.getAllByRole('button').map((button) => button.textContent.replace(/\s+/g, ' ').trim())
	).toEqual(['🇺🇸 English', '🇩🇪 Deutsch', '🇪🇸 Español']);
	expect(option('English')).toHaveAttribute('aria-current', 'true');
	expect(option('Deutsch')).toHaveAttribute('aria-current', 'false');
});

test('switches the language live and closes when picking a locale', async () => {
	const user = userEvent.setup();
	render(LanguageSelector);

	await user.click(trigger());
	await user.click(option('Español'));

	expect(dropdown()).not.toHaveAttribute('open');
	// no reload: the trigger and the menu re-render in Spanish
	expect(trigger()).toHaveTextContent('Idioma');
	await user.click(trigger());
	expect(screen.getByText('Seleccionar idioma')).toBeVisible();
	expect(option('Español')).toHaveAttribute('aria-current', 'true');
	expect(option('English')).toHaveAttribute('aria-current', 'false');
});
