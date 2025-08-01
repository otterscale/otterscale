<script lang="ts">
	import { setContext } from 'svelte';
	import { ModeWatcher } from 'mode-watcher';
	import { createConnectTransport } from '@connectrpc/connect-web';
	import { addCollection } from '@iconify/svelte';
	import logosIcons from '@iconify-json/logos/icons.json';
	import phIcons from '@iconify-json/ph/icons.json';
	import streamlineLogosIcons from '@iconify-json/streamline-logos/icons.json';
	import { env } from '$env/dynamic/public';
	import { Toaster } from '$lib/components/ui/sonner/index.js';

	import '../app.css';

	let { children } = $props();

	const transport = createConnectTransport({
		baseUrl: env.PUBLIC_API_URL
	});

	setContext('transport', transport);

	addCollection(logosIcons);
	addCollection(phIcons);
	addCollection(streamlineLogosIcons);
</script>

<ModeWatcher />
<Toaster richColors />

{@render children()}
