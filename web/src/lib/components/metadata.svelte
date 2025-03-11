<script lang="ts">
	import { page } from '$app/state';
	import { siteConfig } from '$lib/config/site';
	import { i18n } from '$lib/i18n';
	import { featureTitle } from './features';

	const title = $derived.by(() => {
		const canonicalPath = i18n.route(page.url.pathname);
		if (canonicalPath === '/') {
			return `${siteConfig.title} | ${siteConfig.description}`;
		}
		return `${featureTitle(canonicalPath)} | ${siteConfig.title}`;
	});
</script>

<svelte:head>
	<title>{title}</title>
</svelte:head>
