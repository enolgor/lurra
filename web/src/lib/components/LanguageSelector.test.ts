import { render, screen, waitFor } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { afterEach, expect, test } from 'vitest';
import '$lib/i18n.svelte';
import { setLocale } from '../../paraglide/runtime.js';
import LanguageSelector from './LanguageSelector.svelte';

// The menu stays `visibility: hidden` until Floating UI measures its position, which never
// happens in jsdom (no layout). So queries inside it need { hidden: true }, and items are found
// by text: hidden elements have no accessible name.
const inMenu = { hidden: true };
const item = (text: string) => screen.getByText(text).closest('[role="menuitemradio"]');

afterEach(() => {
	void setLocale('en');
	localStorage.clear();
});

test('lists every locale in its own language, with the current one checked', async () => {
	const user = userEvent.setup();
	render(LanguageSelector);

	await user.click(screen.getByRole('button', { name: 'Language' }));

	expect(screen.getByText('Select language')).toBeInTheDocument();
	const items = screen.getAllByRole('menuitemradio', inMenu);
	expect(items.map((item) => item.textContent?.trim())).toEqual([
		'🇺🇸 English',
		'🇩🇪 Deutsch',
		'🇪🇸 Español'
	]);
	expect(item('English')).toHaveAttribute('aria-checked', 'true');
	expect(item('Deutsch')).toHaveAttribute('aria-checked', 'false');
});

test('switches the language live when picking a locale', async () => {
	const user = userEvent.setup();
	render(LanguageSelector);

	await user.click(screen.getByRole('button', { name: 'Language' }));
	await user.click(item('Español')!);

	// no reload: the trigger re-renders in Spanish
	expect(await screen.findByRole('button', { name: 'Idioma' })).toBeInTheDocument();
	await user.click(screen.getByRole('button', { name: 'Idioma' }));
	expect(screen.getByText('Seleccionar idioma')).toBeInTheDocument();
	expect(item('Español')).toHaveAttribute('aria-checked', 'true');
	expect(item('English')).toHaveAttribute('aria-checked', 'false');
});

test('closes the menu after picking a locale', async () => {
	const user = userEvent.setup();
	render(LanguageSelector);

	const trigger = screen.getByRole('button', { name: 'Language' });
	await user.click(trigger);
	await user.click(item('Español')!);

	expect(trigger).toHaveAttribute('aria-expanded', 'false');
	// removed after its exit animation
	await waitFor(() => expect(screen.queryByRole('menu', inMenu)).not.toBeInTheDocument());
});
