<script lang="ts">
	import { page } from '$app/state';
	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { Dashboard } from '$lib/components/applications/dashboard';
	import Loading from '$lib/components/custom/loading/report.svelte';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onMount } from 'svelte';

	import * as Tabs from '$lib/components/ui/tabs';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import TrendingUpIcon from '@lucide/svelte/icons/trending-up';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.applications(page.params.scope) });

	//
	const chartData1 = [
		{ date: new Date('2024-01-01'), desktop: 186, mobile: 80 },
		{ date: new Date('2024-02-01'), desktop: 305, mobile: 200 },
		{ date: new Date('2024-03-01'), desktop: 237, mobile: 120 },
		{ date: new Date('2024-04-01'), desktop: 73, mobile: 190 },
		{ date: new Date('2024-05-01'), desktop: 209, mobile: 130 },
		{ date: new Date('2024-06-01'), desktop: 214, mobile: 140 }
	];
	const chartConfig1 = {
		desktop: { label: 'Desktop', color: 'var(--chart-1)' },
		mobile: { label: 'Mobile', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	//
	const chartData2 = [
		{ date: new Date('2024-01-01'), desktop: 186, mobile: 80 },
		{ date: new Date('2024-02-01'), desktop: 305, mobile: 200 },
		{ date: new Date('2024-03-01'), desktop: 237, mobile: 120 },
		{ date: new Date('2024-04-01'), desktop: 73, mobile: 190 },
		{ date: new Date('2024-05-01'), desktop: 209, mobile: 130 },
		{ date: new Date('2024-06-01'), desktop: 214, mobile: 140 }
	];
	const chartConfig2 = {
		desktop: { label: 'Desktop', color: 'var(--chart-1)' },
		mobile: { label: 'Mobile', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	//
	const chartData3 = [
		{ date: new Date('2024-01-01'), desktop: 186, mobile: 80 },
		{ date: new Date('2024-02-01'), desktop: 305, mobile: 200 },
		{ date: new Date('2024-03-01'), desktop: 237, mobile: 120 },
		{ date: new Date('2024-04-01'), desktop: 73, mobile: 190 },
		{ date: new Date('2024-05-01'), desktop: 209, mobile: 130 },
		{ date: new Date('2024-06-01'), desktop: 214, mobile: 140 }
	];
	const chartConfig3 = {
		desktop: { label: 'Desktop', color: 'var(--chart-1)' },
		mobile: { label: 'Mobile', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	//
	const chartData4 = [
		{ date: new Date('2024-01-01'), desktop: 186, mobile: 80 },
		{ date: new Date('2024-02-01'), desktop: 305, mobile: 200 },
		{ date: new Date('2024-03-01'), desktop: 237, mobile: 120 },
		{ date: new Date('2024-04-01'), desktop: 73, mobile: 190 },
		{ date: new Date('2024-05-01'), desktop: 209, mobile: 130 },
		{ date: new Date('2024-06-01'), desktop: 214, mobile: 140 }
	];
	const chartConfig4 = {
		desktop: { label: 'Desktop', color: 'var(--chart-1)' },
		mobile: { label: 'Mobile', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	//
	const transport: Transport = getContext('transport');
	const environmentService = createClient(EnvironmentService, transport);

	let prometheusDriver: PrometheusDriver | null = null;
	let mounted = false;

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
			<Card.Root class="col-span-2 gap-2">
				<Card.Header>
					<Card.Title>Control Plane</Card.Title>
					<Card.Description>xxx</Card.Description>
				</Card.Header>
				<Card.Content>1 / 1</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-2 gap-2">
				<Card.Header>
					<Card.Title>Worker</Card.Title>
					<Card.Description>xxx</Card.Description>
				</Card.Header>
				<Card.Content>2 / 3</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-3 row-span-2 gap-2">
				<Card.Header>
					<Card.Title>CPU Usage</Card.Title>
					<Card.Description>不區分 namespace</Card.Description>
					<Card.Action class="text-muted-foreground flex flex-col gap-0.5 text-sm">
						<span>Requests: 5%</span>
						<span>Limits: 0.946%</span>
					</Card.Action>
				</Card.Header>

				<Card.Content>
					<Chart.Container config={chartConfig1}>
						<AreaChart
							data={chartData1}
							x="date"
							xScale={scaleUtc()}
							yPadding={[0, 25]}
							series={[
								{
									key: 'mobile',
									label: 'Mobile',
									color: 'var(--color-mobile)'
								}
							]}
							seriesLayout="stack"
							props={{
								area: {
									curve: curveNatural,
									'fill-opacity': 0.4,
									line: { class: 'stroke-1' },
									motion: 'tween'
								},
								xAxis: {
									format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' })
								},
								yAxis: { format: () => '' }
							}}
						>
							{#snippet tooltip()}
								<Chart.Tooltip
									indicator="dot"
									labelFormatter={(v: Date) => {
										return v.toLocaleDateString('en-US', {
											month: 'long'
										});
									}}
								/>
							{/snippet}
							{#snippet marks({ series, getAreaProps })}
								{#each series as s, i (s.key)}
									<LinearGradient
										stops={[s.color ?? '', 'color-mix(in lch, ' + s.color + ' 10%, transparent)']}
										vertical
									>
										{#snippet children({ gradient })}
											<Area {...getAreaProps(s, i)} fill={gradient} />
										{/snippet}
									</LinearGradient>
								{/each}
							{/snippet}
						</AreaChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-3 row-span-2 gap-2">
				<Card.Header>
					<Card.Title>Memory Usage (w/o cache)</Card.Title>
					<Card.Description>不區分 namespace</Card.Description>
					<Card.Action class="text-muted-foreground flex flex-col gap-0.5 text-sm">
						<span>Requests: 32.3%</span>
						<span>Limits: 1.91%</span>
					</Card.Action>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={chartConfig2}>
						<AreaChart
							data={chartData2}
							x="date"
							xScale={scaleUtc()}
							yPadding={[0, 25]}
							series={[
								{
									key: 'desktop',
									label: 'Desktop',
									color: 'var(--color-desktop)'
								}
							]}
							seriesLayout="stack"
							props={{
								area: {
									curve: curveNatural,
									'fill-opacity': 0.4,
									line: { class: 'stroke-1' },
									motion: 'tween'
								},
								xAxis: {
									format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' })
								},
								yAxis: { format: () => '' }
							}}
						>
							{#snippet tooltip()}
								<Chart.Tooltip
									indicator="dot"
									labelFormatter={(v: Date) => {
										return v.toLocaleDateString('en-US', {
											month: 'long'
										});
									}}
								/>
							{/snippet}
							{#snippet marks({ series, getAreaProps })}
								{#each series as s, i (s.key)}
									<LinearGradient
										stops={[s.color ?? '', 'color-mix(in lch, ' + s.color + ' 10%, transparent)']}
										vertical
									>
										{#snippet children({ gradient })}
											<Area {...getAreaProps(s, i)} fill={gradient} />
										{/snippet}
									</LinearGradient>
								{/each}
							{/snippet}
						</AreaChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-2 col-start-1 gap-2">
				<Card.Header>
					<Card.Title>Running Pods</Card.Title>
					<Card.Description>xxx</Card.Description>
				</Card.Header>
				<Card.Content>52</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-2 gap-2">
				<Card.Header>
					<Card.Title>Running Containers</Card.Title>
					<Card.Description>xxx</Card.Description>
				</Card.Header>
				<Card.Content>136</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-3 col-start-5 gap-2">
				<Card.Header>
					<Card.Title>Network Bandwidth</Card.Title>
					<Card.Description>Recieve & Transmit</Card.Description>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={chartConfig3}>
						<AreaChart
							data={chartData3}
							x="date"
							xScale={scaleUtc()}
							yPadding={[0, 25]}
							series={[
								{
									key: 'mobile',
									label: 'Mobile',
									color: 'var(--color-mobile)'
								},
								{
									key: 'desktop',
									label: 'Desktop',
									color: 'var(--color-desktop)'
								}
							]}
							seriesLayout="stack"
							props={{
								area: {
									curve: curveNatural,
									'fill-opacity': 0.4,
									line: { class: 'stroke-1' },
									motion: 'tween'
								},
								xAxis: {
									format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' })
								},
								yAxis: { format: () => '' }
							}}
						>
							{#snippet tooltip()}
								<Chart.Tooltip
									indicator="dot"
									labelFormatter={(v: Date) => {
										return v.toLocaleDateString('en-US', {
											month: 'long'
										});
									}}
								/>
							{/snippet}
							{#snippet marks({ series, getAreaProps })}
								{#each series as s, i (s.key)}
									<LinearGradient
										stops={[s.color ?? '', 'color-mix(in lch, ' + s.color + ' 10%, transparent)']}
										vertical
									>
										{#snippet children({ gradient })}
											<Area {...getAreaProps(s, i)} fill={gradient} />
										{/snippet}
									</LinearGradient>
								{/each}
							{/snippet}
						</AreaChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-3 gap-2">
				<Card.Header>
					<Card.Title>Storage ThroughPut</Card.Title>
					<Card.Description>Read & Write</Card.Description>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={chartConfig4}>
						<AreaChart
							data={chartData4}
							x="date"
							xScale={scaleUtc()}
							yPadding={[0, 25]}
							series={[
								{
									key: 'mobile',
									label: 'Mobile',
									color: 'var(--color-mobile)'
								},
								{
									key: 'desktop',
									label: 'Desktop',
									color: 'var(--color-desktop)'
								}
							]}
							seriesLayout="stack"
							props={{
								area: {
									curve: curveNatural,
									'fill-opacity': 0.4,
									line: { class: 'stroke-1' },
									motion: 'tween'
								},
								xAxis: {
									format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' })
								},
								yAxis: { format: () => '' }
							}}
						>
							{#snippet tooltip()}
								<Chart.Tooltip
									indicator="dot"
									labelFormatter={(v: Date) => {
										return v.toLocaleDateString('en-US', {
											month: 'long'
										});
									}}
								/>
							{/snippet}
							{#snippet marks({ series, getAreaProps })}
								{#each series as s, i (s.key)}
									<LinearGradient
										stops={[s.color ?? '', 'color-mix(in lch, ' + s.color + ' 10%, transparent)']}
										vertical
									>
										{#snippet children({ gradient })}
											<Area {...getAreaProps(s, i)} fill={gradient} />
										{/snippet}
									</LinearGradient>
								{/each}
							{/snippet}
						</AreaChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>
		</Tabs.Content>
		<Tabs.Content value="analytics">
			{#if mounted && prometheusDriver && $activeScope}
				<Dashboard client={prometheusDriver} scope={$activeScope} />
			{:else}
				<Loading />
			{/if}
		</Tabs.Content>
	</Tabs.Root>
</div>
