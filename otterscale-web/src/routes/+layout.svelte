<script lang="ts">
	import { setContext, onMount } from 'svelte';
	import { ModeWatcher } from 'mode-watcher';
	import { createConnectTransport } from '@connectrpc/connect-web';
	import { createClient } from '@connectrpc/connect';
	import { addCollection } from '@iconify/svelte';
	import logosIcons from '@iconify-json/logos/icons.json';
	import phIcons from '@iconify-json/ph/icons.json';
	import streamlineLogosIcons from '@iconify-json/streamline-logos/icons.json';
	import { env } from '$env/dynamic/public';
	import favicon from '$lib/assets/favicon.svg';
	import { Toaster } from '$lib/components/ui/sonner/index.js';
	import { PrometheusDriver } from 'prometheus-query';
	import { writable, type Writable } from 'svelte/store';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import '../app.css';

	let { children } = $props();

	const transport = createConnectTransport({
		baseUrl: env.PUBLIC_API_URL
	});

	setContext('transport', transport);

	// // 創建 PrometheusDriver 初始化 promise
	// const prometheusDriver = (async () => {
	// 	try {
	// 		const environmentService = createClient(EnvironmentService, transport);
	// 		const response = await environmentService.getPrometheus({});
	// 		return new PrometheusDriver({
	// 			endpoint: response.endpoint,
	// 			baseURL: response.baseUrl
	// 		});
	// 	} catch (error) {
	// 		console.error('Error initializing PrometheusDriver:', error);
	// 		throw error;
	// 	}
	// })();

	// setContext('prometheusDriver', prometheusDriver);
	addCollection(logosIcons);
	addCollection(phIcons);
	addCollection(streamlineLogosIcons);
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<ModeWatcher />
<Toaster richColors />

{@render children()}
