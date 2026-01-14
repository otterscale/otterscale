<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';

	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { Reloader } from '$lib/components/custom/reloader';
	import { Overview } from '$lib/components/scopes/overview/index';
	import { m } from '$lib/paraglide/messages';

	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const environmentService = createClient(EnvironmentService, transport);

	let isReloading = $state(true);
	let prometheusDriver = $state<PrometheusDriver | null>(null);

	onMount(async () => {
		try {
			const response = await environmentService.getPrometheus({});
			prometheusDriver = new PrometheusDriver({
				endpoint: `${env.PUBLIC_API_URL}/prometheus`,
				baseURL: response.baseUrl
			});
		} catch (error) {
			console.error('Failed to initialize Prometheus driver:', error);
		}
	});

	onDestroy(() => {
		isReloading = false;
	});
</script>

<main class="space-y-4 py-4">
	{#if prometheusDriver}
		<div class="mx-auto grid w-full gap-6">
			<div class="grid gap-1">
				<h1 class="text-2xl font-bold tracking-tight md:text-3xl">{m.k8s_overview_title()}</h1>
				<p class="text-muted-foreground">
					{m.k8s_overview_description()}
				</p>
			</div>

			<div class="flex justify-end">
				<Reloader bind:checked={isReloading} />
			</div>
			<Overview {prometheusDriver} {scope} bind:isReloading />
		</div>
	{/if}
</main>
