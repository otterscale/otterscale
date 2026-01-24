<script lang="ts">
	import '../app.css';

	import type { Interceptor } from '@connectrpc/connect';
	import { createConnectTransport } from '@connectrpc/connect-web';
	import { addCollection } from '@iconify/svelte';
	import logos from '@iconify-json/logos/icons.json';
	import ph from '@iconify-json/ph/icons.json';
	import simpleIcons from '@iconify-json/simple-icons/icons.json';
	import { ModeWatcher } from 'mode-watcher';
	import { setContext } from 'svelte';

	import { dev } from '$app/environment';
	import { Toaster } from '$lib/components/ui/sonner';
	import { Spinner } from '$lib/components/ui/spinner';

	let { children } = $props();

	const proxyHeaderInterceptor: Interceptor = (next) => async (req) => {
		req.header.set('x-proxy-target', 'api');
		return await next(req);
	};

	const transport = createConnectTransport({
		baseUrl: '/',
		useBinaryFormat: !dev,
		interceptors: [proxyHeaderInterceptor],
		fetch
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
<Toaster expand richColors {loadingIcon} />

<div class="app">
	<main>
		{@render children()}
	</main>
</div>
