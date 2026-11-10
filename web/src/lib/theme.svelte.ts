// Light/dark theme choice. daisyUI themes the page from data-theme on <html>; without it, its
// `dark` theme follows the system (--prefersdark), which is what "system" means here.
// public/theme.js applies the saved choice before the first paint; keep its key in sync.

export const themes = ['light', 'dark', 'system'] as const;
export type Theme = (typeof themes)[number];

export const themeStorageKey = 'theme';

function isTheme(value: unknown): value is Theme {
	return themes.includes(value as Theme);
}

/** The saved choice, else "system". */
export function detectTheme(): Theme {
	let saved: string | null = null;
	try {
		saved = localStorage.getItem(themeStorageKey);
	} catch {
		// storage blocked (private mode, site data disabled): follow the system
	}
	return isTheme(saved) ? saved : 'system';
}

const initial = detectTheme();
let current = $state(initial);

/** The current choice; reactive in markup and $derived. */
export function getTheme(): Theme {
	return current;
}

function applyTheme(theme: Theme) {
	if (theme === 'system') document.documentElement.removeAttribute('data-theme');
	else document.documentElement.dataset.theme = theme;
}

/** Applies the theme and remembers it; "system" forgets the saved choice. */
export function setTheme(theme: Theme) {
	current = theme;
	applyTheme(theme);
	try {
		if (theme === 'system') localStorage.removeItem(themeStorageKey);
		else localStorage.setItem(themeStorageKey, theme);
	} catch {
		// storage blocked: the choice lasts until the page is reloaded
	}
}

applyTheme(initial);
