<script lang="ts">
	import '../app.css';

	import { createConnectTransport } from '@connectrpc/connect-web';
	import { addCollection } from '@iconify/svelte';
	import logos from '@iconify-json/logos/icons.json';
	import ph from '@iconify-json/ph/icons.json';
	import simpleIcons from '@iconify-json/simple-icons/icons.json';
	import { ModeWatcher } from 'mode-watcher';
	import { setContext } from 'svelte';

	import { env } from '$env/dynamic/public';
	import { Toaster } from '$lib/components/ui/sonner';
	import { Spinner } from '$lib/components/ui/spinner';

	let { children } = $props();

	const transport = createConnectTransport({
		baseUrl: env.PUBLIC_API_URL || ''
	});

	setContext('transport', transport);

	addCollection(logos);
	addCollection(ph);
	addCollection(simpleIcons);
</script>

{#snippet loadingIcon()}
	<Spinner />
{/snippet}

<ModeWatcher />
<Toaster closeButton expand richColors {loadingIcon} />

<div class="app">
	<main>
		{@render children()}
	</main>
</div>
