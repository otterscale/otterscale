<script lang="ts">
	import { setContext, onMount } from 'svelte';
	import { ModeWatcher } from 'mode-watcher';
	import { createConnectTransport } from '@connectrpc/connect-web';
	import { addCollection } from '@iconify/svelte';
	import logos from '@iconify-json/logos/icons.json';
	import ph from '@iconify-json/ph/icons.json';
	import simpleIcons from '@iconify-json/simple-icons/icons.json';
	import streamlineLogos from '@iconify-json/streamline-logos/icons.json';
	import { env } from '$env/dynamic/public';
	import favicon from '$lib/assets/favicon.svg';
	import { Toaster } from '$lib/components/ui/sonner/index.js';
	import '../app.css';

	let { children } = $props();

	const transport = createConnectTransport({
		baseUrl: env.PUBLIC_API_URL,
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

<ModeWatcher />
<Toaster richColors />

{@render children()}
