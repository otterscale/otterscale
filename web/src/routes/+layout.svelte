<script lang="ts">
	import { createConnectTransport } from '@connectrpc/connect-web';
	import { addCollection } from '@iconify/svelte';
    import Icon from '@iconify/svelte';
	import logos from '@iconify-json/logos/icons.json';
	import ph from '@iconify-json/ph/icons.json';
	import simpleIcons from '@iconify-json/simple-icons/icons.json';
	import streamlineLogos from '@iconify-json/streamline-logos/icons.json';
	import { ModeWatcher } from 'mode-watcher';
	import { setContext } from 'svelte';

	import { env } from '$env/dynamic/public';
	import favicon from '$lib/assets/favicon.svg';
	import { Toaster } from '$lib/components/ui/sonner/index.js';
	import '../app.css';

	let { children } = $props();

	const transport = createConnectTransport({
		baseUrl: env.PUBLIC_API_URL || '',
	});

	setContext('transport', transport);

	addCollection(logos);
	addCollection(ph);
	addCollection(simpleIcons);
	addCollection(streamlineLogos);
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{#snippet loadingIcon()}
    <Icon icon="ph:spinner" class="animate-spin" />
{/snippet}

<ModeWatcher />
<Toaster richColors {loadingIcon}/>

{@render children()}
