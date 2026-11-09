// Live locale switching: the current locale lives in $state, and Paraglide's getLocale() reads
// it, so every m.key() in markup or $derived re-renders on setLocale(), without a page reload.
// Import this module before anything calls m.key() (main.ts does).
import {
	baseLocale,
	extractLocaleFromNavigator,
	getTextDirection,
	isLocale,
	localStorageKey,
	overwriteGetLocale,
	overwriteSetLocale,
	type Locale
} from '$lib/paraglide/runtime.js';

/** The saved choice, else the browser's preferred languages, else the base locale (en). */
export function detectLocale(): Locale {
	let saved: string | null = null;
	try {
		saved = localStorage.getItem(localStorageKey);
	} catch {
		// storage blocked (private mode, site data disabled): fall through to the browser language
	}
	if (isLocale(saved)) return saved;
	return extractLocaleFromNavigator() ?? baseLocale;
}

const initial = detectLocale();
let current = $state(initial);

/** Sets <html lang> and dir, for screen readers, hyphenation and right-to-left languages. */
function updateDocument(locale: Locale) {
	document.documentElement.lang = locale;
	document.documentElement.dir = getTextDirection(locale);
}

overwriteGetLocale(() => current);
overwriteSetLocale((locale) => {
	current = locale;
	updateDocument(locale);
	try {
		localStorage.setItem(localStorageKey, locale);
	} catch {
		// storage blocked: the choice lasts until the page is reloaded
	}
});
updateDocument(initial);
