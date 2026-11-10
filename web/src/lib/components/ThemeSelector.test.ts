import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { afterEach, expect, test } from 'vitest';
import { setLocale } from '$lib/paraglide/runtime.js';
import { setTheme } from '$lib/theme.svelte';
import ThemeSelector from './ThemeSelector.svelte';
import { trigger } from '$lib/test/dropdown';

const triggerIcon = (name: string) => trigger().querySelector(`.lucide-${name}`);
const option = (name: string) => screen.getByRole('button', { name });

afterEach(() => {
	setTheme('system');
});

test('a "Theme" button showing the current choice, system by default', () => {
	render(ThemeSelector);

	expect(trigger()).toHaveTextContent('Theme');
	expect(triggerIcon('monitor')).toBeInTheDocument();
});

test('lists light, dark and system, with the current one marked', async () => {
	const user = userEvent.setup();
	render(ThemeSelector);

	await user.click(trigger());

	expect(screen.getByText('Select theme')).toBeVisible();
	expect(screen.getAllByRole('button').map((button) => button.textContent.trim())).toEqual([
		'Light',
		'Dark',
		'System'
	]);
	expect(option('System')).toHaveAttribute('aria-current', 'true');
	expect(option('Dark')).toHaveAttribute('aria-current', 'false');
});

test('picking a theme applies it and updates the button and the mark', async () => {
	const user = userEvent.setup();
	render(ThemeSelector);

	await user.click(trigger());
	await user.click(option('Dark'));

	expect(document.documentElement.dataset.theme).toBe('dark');
	expect(triggerIcon('moon')).toBeInTheDocument();
	await user.click(trigger());
	expect(option('Dark')).toHaveAttribute('aria-current', 'true');
	expect(option('System')).toHaveAttribute('aria-current', 'false');

	await user.click(option('Light'));
	expect(triggerIcon('sun')).toBeInTheDocument();
});

test('labels follow the language live', async () => {
	const user = userEvent.setup();
	render(ThemeSelector);

	void setLocale('es');
	await user.click(trigger());

	expect(await screen.findByRole('button', { name: 'Oscuro' })).toBeVisible();
	expect(trigger()).toHaveTextContent('Tema');
});
