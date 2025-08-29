<script lang="ts">
	import { page } from '$app/state';
	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import Loading from '$lib/components/custom/loading/report.svelte';
	import { ReloadManager, Reloader } from '$lib/components/custom/reloader';
	import { DashboardAnalytics } from '$lib/components/storage/dashboard';
	import { DashboardOverview } from '$lib/components/storage/dashboard-overview';
	import * as Tabs from '$lib/components/ui/tabs';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.storage(page.params.scope) });

	const transport: Transport = getContext('transport');
	const environmentService = createClient(EnvironmentService, transport);

	let prometheusDriver = $state<PrometheusDriver | null>(null);
	let mounted = $state(false);

	async function initializePrometheusDriver(): Promise<PrometheusDriver> {
		const response = await environmentService.getPrometheus({});
		return new PrometheusDriver({
			endpoint: `${env.PUBLIC_API_URL}/prometheus`,
			baseURL: response.baseUrl
		});
	}
	let counter = $state(0);
	const reloadManager = new ReloadManager(() => {
		counter = counter + 1;
	});

	onMount(async () => {
		try {
			prometheusDriver = await initializePrometheusDriver();
		} catch (error) {
			console.error('Failed to initialize Prometheus driver:', error);
		} finally {
			mounted = true;
		}
		// reloadManager.start();
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<div class="mx-auto grid w-full gap-6">
	<div class="grid gap-1">
		<h1 class="text-2xl font-bold tracking-tight md:text-3xl">Dashboard</h1>
		<p class="text-muted-foreground">description</p>
	</div>

	<Tabs.Root value="overview">
		<div class="flex justify-between gap-2">
			<Tabs.List>
				<Tabs.Trigger value="overview">{m.overview()}</Tabs.Trigger>
				<Tabs.Trigger value="analytics">{m.analytics()}</Tabs.Trigger>
			</Tabs.List>
			<Reloader {reloadManager} />
		</div>
		<Tabs.Content
			value="overview"
			class="grid auto-rows-auto grid-cols-2 gap-5 pt-4 md:grid-cols-4 lg:grid-cols-10"
		>
			{#if mounted && prometheusDriver && $activeScope}
				{#key counter}
					<DashboardOverview client={prometheusDriver} scope={$activeScope} />
				{/key}
			{:else}
				<Loading />
			{/if}
		</Tabs.Content>
		<Tabs.Content value="analytics">
			{#if mounted && prometheusDriver && $activeScope}
				{#key counter}
					<DashboardAnalytics client={prometheusDriver} scope={$activeScope} />
				{/key}
			{:else}
				<Loading />
			{/if}
		</Tabs.Content>
	</Tabs.Root>
</div>
