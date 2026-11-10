import { defined } from './defined';

/** The `<details>` of the Dropdown on the page; it has `open` while the menu shows. */
export const dropdown = () => defined(document.querySelector('details'));
/** The Dropdown's button (its `<summary>`), which opens and closes the menu. */
export const trigger = () => defined(document.querySelector('summary'));
