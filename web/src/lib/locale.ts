// Display helpers for locales, derived from Intl so new locales need no extra table.

/** The locale's name in its own language, capitalized: "English", "Español", "Deutsch". */
export function localeName(locale: string): string {
	const name = String(
		new Intl.DisplayNames([locale], { type: 'language', fallback: 'code' }).of(locale)
	);
	return name.charAt(0).toLocaleUpperCase(locale) + name.slice(1);
}

/**
 * The flag emoji of the locale's most likely country (en → 🇺🇸, es → 🇪🇸), or '' when it has
 * none (eo → 001, the world). Systems without flag emoji (Windows) show the letters ("ES").
 */
export function localeFlag(locale: string): string {
	const region = String(new Intl.Locale(locale).maximize().region);
	if (!/^[A-Z]{2}$/.test(region)) return '';
	return String.fromCodePoint(...[...region].map((c) => 0x1f1e6 + c.charCodeAt(0) - 65));
}
