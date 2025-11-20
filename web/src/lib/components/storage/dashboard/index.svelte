<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';

	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import Reloader from '$lib/components/custom/reloader/reloader.svelte';
	import { Overview } from '$lib/components/storage/dashboard/overview';
	import * as Tabs from '$lib/components/ui/tabs';
	import { m } from '$lib/paraglide/messages';

	import ExtensionsAlert from './extensions-alert.svelte';

	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const environmentService = createClient(EnvironmentService, transport);

	let prometheusDriver = $state<PrometheusDriver | null>(null);

	let isReloading = $state(true);
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
	<ExtensionsAlert {scope} />
	{#if prometheusDriver}
		<div class="mx-auto grid w-full gap-6">
			<div class="grid gap-1">
				<h1 class="text-2xl font-bold tracking-tight md:text-3xl">{m.storage()}</h1>
				<p class="text-muted-foreground">
					{m.storage_dashboard_description()}
				</p>
			</div>
			<Tabs.Root value="overview">
				<div class="flex justify-between gap-2">
					<Tabs.List>
						<Tabs.Trigger value="overview">{m.overview()}</Tabs.Trigger>
						<Tabs.Trigger value="analytics" disabled>{m.analytics()}</Tabs.Trigger>
					</Tabs.List>
					<Reloader bind:checked={isReloading} />
				</div>
				<Tabs.Content value="overview">
					<Overview client={prometheusDriver} {scope} bind:isReloading />
				</Tabs.Content>
				<Tabs.Content value="analytics">
					<!-- <Analytics client={prometheusDriver} {scope} /> -->
				</Tabs.Content>
			</Tabs.Root>
		</div>
	{/if}
</main>
