import { mount } from 'svelte';
import './lib/i18n.svelte'; // before App: sets up live locale switching
import App from './App.svelte';

mount(App, { target: document.querySelector('#app')! });
