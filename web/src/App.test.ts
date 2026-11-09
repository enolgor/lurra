import { render, screen, within } from '@testing-library/svelte';
import { afterEach, expect, test } from 'vitest';
import '$lib/i18n.svelte';
import { setLocale } from '$lib/paraglide/runtime.js';
import App from './App.svelte';

afterEach(() => {
	void setLocale('en');
	localStorage.clear();
});

test('has a header with the app title and settings, a main area and a footer', async () => {
	history.replaceState(null, '', '/');
	render(App);

	const header = screen.getByRole('banner');
	expect(within(header).getByRole('link', { name: 'Home' })).toHaveAttribute('href', '/');
	// the two menu buttons (their menus repeat the names as titles)
	expect(within(header).getByText('Theme', { selector: 'summary span' })).toBeInTheDocument();
	expect(within(header).getByText('Language', { selector: 'summary span' })).toBeInTheDocument();

	const main = screen.getByRole('main');
	expect(main).toHaveAttribute('id', 'main');
	expect(await within(main).findByRole('heading', { level: 1, name: 'Todos' })).toBeVisible();

	const footer = screen.getByRole('contentinfo');
	expect(footer).toHaveTextContent('An example app built with the lurra project template');
	expect(within(footer).getByRole('link', { name: 'lurra project template' })).toHaveAttribute(
		'href',
		'https://github.com/enolgor/lurra'
	);
});

test('starts with a skip link to the main area', () => {
	render(App);

	const skip = screen.getAllByRole('link')[0];
	expect(skip).toHaveTextContent('Skip to content');
	expect(skip).toHaveAttribute('href', '#main');
});

test('the browser tab title is the app title, in the current language', async () => {
	render(App);
	expect(document.title).toBe('Todos');

	void setLocale('es');
	await screen.findByRole('link', { name: 'Inicio' });
	expect(document.title).toBe('Tareas');
});

test.each(['/', '/todos/1'])('renders the page for %s inside main', async (path) => {
	history.replaceState(null, '', path);
	render(App);

	expect(await within(screen.getByRole('main')).findByRole('heading', { level: 1 })).toBeVisible();
});
