<script lang="ts">
	import { setContext } from 'svelte';
	import { ModeWatcher } from 'mode-watcher';
	import { createConnectTransport } from '@connectrpc/connect-web';
	import { PUBLIC_API_URL } from '$env/static/public';
	import { Toaster } from '$lib/components/ui/sonner/index.js';
	import '../app.css';

	let { children } = $props();

	const transport = createConnectTransport({
		baseUrl: PUBLIC_API_URL
	});

	setContext('transport', transport);
</script>

<ModeWatcher />
<Toaster />

{@render children()}
