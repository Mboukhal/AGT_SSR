// @ts-check
import { defineConfig } from 'astro/config';

import svelte from '@astrojs/svelte';

import tailwindcss from '@tailwindcss/vite';

// https://astro.build/config
export default defineConfig({
  server: {
    port: 1337,
  },

  integrations: [svelte()],

  vite: {
    envDir: "../",
    clearScreen: false,
    plugins: [tailwindcss()],
  },

});