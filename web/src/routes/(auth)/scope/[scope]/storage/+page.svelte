<script lang="ts">
	import { page } from '$app/state';
	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import Loading from '$lib/components/custom/loading/report.svelte';
	import { DashboardOverview } from '$lib/components/storage/dashboard-overview';
	import { DashboardAnalytics } from '$lib/components/storage/dashboard';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import * as Tabs from '$lib/components/ui/tabs';
	import { getLocale } from '$lib/paraglide/runtime';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import ActivityIcon from '@lucide/svelte/icons/activity';
	import { scaleUtc } from 'd3-scale';
	import { curveLinear, curveStep } from 'd3-shape';
	import {
		ArcChart,
		AreaChart,
		BarChart,
		type ChartContextValue,
		Highlight,
		LineChart,
		PieChart,
		Text
	} from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onMount } from 'svelte';
	import { cubicInOut } from 'svelte/easing';

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

	onMount(async () => {
		try {
			prometheusDriver = await initializePrometheusDriver();
		} catch (error) {
			console.error('Failed to initialize Prometheus driver:', error);
		} finally {
			mounted = true;
		}
	});
</script>

<div class="mx-auto grid w-full gap-6">
	<div class="grid gap-1">
		<h1 class="text-2xl font-bold tracking-tight md:text-3xl">Dashboard</h1>
		<p class="text-muted-foreground">description</p>
	</div>

	<Tabs.Root value="overview">
		<Tabs.List>
			<Tabs.Trigger value="overview">Overview</Tabs.Trigger>
			<Tabs.Trigger value="analytics">Analytics</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content
			value="overview"
			class="grid auto-rows-auto grid-cols-2 gap-5 pt-4 md:grid-cols-4 lg:grid-cols-10"
		>
			{#if mounted && prometheusDriver && $activeScope}
				<DashboardOverview client={prometheusDriver} scope={$activeScope} />
			{:else}
				<Loading />
			{/if}
		</Tabs.Content>
		<Tabs.Content value="analytics">
			<!-- {#if mounted && prometheusDriver && $activeScope}
				<DashboardAnalytics client={prometheusDriver} scope={$activeScope} />
			{:else}
				<Loading />
			{/if} -->
		</Tabs.Content>
	</Tabs.Root>
</div>
