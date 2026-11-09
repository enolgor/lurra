// Applies the saved theme before the first paint, so a dark choice never flashes light. A plain
// script in <head> (inline scripts are blocked by the CSP); src/lib/theme.svelte.ts owns the
// choice, uses the same key, and removes it for "system".
try {
	const theme = localStorage.getItem('theme');
	if (theme === 'light' || theme === 'dark') document.documentElement.dataset.theme = theme;
} catch {
	// storage blocked: follow the system
}
