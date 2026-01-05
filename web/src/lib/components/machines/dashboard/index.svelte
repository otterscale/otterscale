<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';

	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { Reloader } from '$lib/components/custom/reloader';
	import { Overview } from '$lib/components/machines/dashboard/overview/index';
	import * as Tabs from '$lib/components/ui/tabs';
	import { m } from '$lib/paraglide/messages';

	import { Dashboard } from './analytics';

	const transport: Transport = getContext('transport');
	const environmentService = createClient(EnvironmentService, transport);

	let isReloading = $state(true);
	let prometheusDriver = $state<PrometheusDriver | null>(null);
	let selectedTab = $state('overview');

	onMount(async () => {
		try {
			const response = await environmentService.getPrometheus({});
			prometheusDriver = new PrometheusDriver({
				endpoint: '/prometheus',
				baseURL: response.baseUrl,
				headers: {
					'x-proxy-target': 'api'
				}
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
				<h1 class="text-2xl font-bold tracking-tight md:text-3xl">{m.dashboard()}</h1>
				<p class="text-muted-foreground">
					{m.machine_dashboard_description()}
				</p>
			</div>

			<Tabs.Root bind:value={selectedTab}>
				<div class="flex justify-between gap-2">
					<Tabs.List>
						<Tabs.Trigger value="overview">{m.overview()}</Tabs.Trigger>
						<Tabs.Trigger value="analytics">{m.analytics()}</Tabs.Trigger>
					</Tabs.List>
					<Reloader bind:checked={isReloading} />
				</div>
				<Tabs.Content value="overview">
					<Overview {selectedTab} {prometheusDriver} bind:isReloading />
				</Tabs.Content>
				<Tabs.Content value="analytics">
					<Dashboard {selectedTab} client={prometheusDriver} />
				</Tabs.Content>
			</Tabs.Root>
		</div>
	{/if}
</main>
