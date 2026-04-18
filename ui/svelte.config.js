import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	compilerOptions: {
		// Force runes mode for the project, except for libraries. Can be removed in svelte 6.
		runes: ({ filename }) => (filename.split(/[/\\]/).includes('node_modules') ? undefined : true)
	},
	adapter: adapter({
		pages: 'build', // path for HTML pages
		assets: 'build', // path for static assets
		fallback: 'index.html', // IMPORTANT
		strict: false // allow dynamic routes,
	}),
	prerender: { entries: ['*'] },
	alias: {
		'@/*': './src/*'
	}
};

export default config;
