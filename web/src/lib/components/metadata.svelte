<script lang="ts">
	import { page } from '$app/state';
	import { siteConfig } from '$lib/config/site';
	import { i18n } from '$lib/i18n';
	import { getFeatureTitle } from '$lib/utils';

	const getPathNodes = (url: URL, language: string): string[] => {
		const pathname =
			language === 'en' ? url.pathname : (url.pathname + '/').replace(`/${language}/`, '/');
		return pathname.split('/').filter((node) => node !== '');
	};

	const title = $derived.by(() => {
		const language = i18n.getLanguageFromUrl(page.url);
		const nodes = getPathNodes(page.url, language);
		if (nodes.length === 0) {
			return `${siteConfig.title} | ${siteConfig.description}`;
		}
		return getFeatureTitle('/' + nodes[0]);
	});
</script>

<svelte:head>
	<title>{title}</title>
</svelte:head>
