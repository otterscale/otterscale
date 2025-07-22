<script lang="ts">
	import { Metadata } from '$lib/components';
	import { Toaster } from '$lib/components/ui/sonner';
	import { i18n } from '$lib/i18n';
	import { createConnectTransport } from '@connectrpc/connect-web';
	import '@fontsource-variable/noto-sans-tc';
	import logosIcons from '@iconify-json/logos/icons.json';
	import phIcons from '@iconify-json/ph/icons.json';
	import { addCollection } from '@iconify/svelte';
	import { ParaglideJS } from '@inlang/paraglide-sveltekit';
	import 'inter-ui/inter-variable.css';
	import { ModeWatcher } from 'mode-watcher';
	import { setContext } from 'svelte';
	import '../app.css';
	import { getContext, onMount } from 'svelte';
	import { createClient } from '@connectrpc/connect';
	import { EnvironmentService, type Prometheus } from '$gen/api/environment/v1/environment_pb';
	import { PrometheusDriver } from 'prometheus-query';
	import { writable, type Writable } from 'svelte/store';

	addCollection(phIcons);
	addCollection(logosIcons);

	let { children } = $props();

	const transport = createConnectTransport({
		baseUrl: import.meta.env.PUBLIC_API_URL
	});
	setContext('transport', transport);

	const prometheusDriver: Writable<PrometheusDriver> = writable({} as PrometheusDriver);
	setContext('prometheusDriver', prometheusDriver);

	onMount(async () => {
		try {
			const environmentService = createClient(EnvironmentService, transport);
			const response = await environmentService.getPrometheus({});
			prometheusDriver.set(
				new PrometheusDriver({
					endpoint: response.endpoint,
					baseURL: response.baseUrl
				})
			);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<Metadata />

<ModeWatcher />
<Toaster closeButton richColors />

<ParaglideJS {i18n}>
	{@render children()}
</ParaglideJS>
