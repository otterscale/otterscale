<script lang="ts">
	const healthData = group(healthRawData, (d) => d.type).get('error') ?? [];
	function metricColor(metric: number) {
		switch (true) {
			case metric > 62:
				return 'fill-red-800';
			case metric > 38:
				return 'fill-yellow-500';
			default:
				return 'fill-green-800';
		}
	}
	function metricBackgroundColor(metric: number) {
		switch (true) {
			case metric > 62:
				return 'fill-red-50';
			case metric > 38:
				return 'fill-yellow-50';
			default:
				return 'fill-green-50';
		}
	}

	const chartHeight = {
		large: 150,
		small: 100
	};
	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	let data = $state(healthData);

	import * as Collapsible from '$lib/components/ui/collapsible';
	import { PageLoading } from '$lib/components/otterscale/ui/index';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Label } from '$lib/components/ui/label';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { cn } from '$lib/utils';
	import { healthRawData } from './dataset';
	import { group } from 'd3-array';
	import { AreaChart, Svg, Group, Arc, Chart, Text } from 'layerchart';
	import { sum } from 'd3-array';
	import { Badge } from '$lib/components/ui/badge';
	import { latencies, storages, inputOutputs, usages, currentUsage } from './dataset';
	import * as Card from '$lib/components/ui/card';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Tabs from '$lib/components/ui/tabs';
	import { Button } from '$lib/components/ui/button';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { writable } from 'svelte/store';
	import PrometheusError from './error/prometheus.svelte';
	import KubernetesError from './error/kubernetes.svelte';
	import CephError from './error/ceph.svelte';
	import Icon from '@iconify/svelte';
	import {
		Nexus,
		type CreateCephRequest,
		type Scope,
		type Machine,
		type Error
	} from '$gen/api/nexus/v1/nexus_pb';

	let { scope }: { scope: Scope } = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const errorsStore = writable<Error[]>([]);
	const errorsLoading = writable(true);
	async function fetchErrors(scopeUuid: string) {
		try {
			const response = await client.verifyEnvironment({ scopeUuid: scopeUuid });
			errorsStore.set(response.errors);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			errorsLoading.set(false);
		}
	}
	async function refreshErrors(scopeUuid: string) {
		while (true) {
			await new Promise((resolve) => setTimeout(resolve, 1000 * 30));

			try {
				const response = await client.verifyEnvironment({ scopeUuid: scopeUuid });
				errorsStore.set(response.errors);
			} catch (error) {
				console.error('Error fetching:', error);
			}
		}
	}

	const currentInformationStore = writable<Error>();
	async function iterateInformations() {
		while (true) {
			for (const error of $errorsStore.filter((e) => Number(e.level) <= 3)) {
				currentInformationStore.set(error);
				await new Promise((resolve) => setTimeout(resolve, 1000 * 1));
			}

			await new Promise((resolve) => setTimeout(resolve, 1000 * 1));
		}
	}

	function isCephError(error: Error) {
		return error.code === 'CEPH_NOT_FOUND';
	}
	function isKubernetesError(error: Error) {
		return error.code === 'KUBERNETES_NOT_FOUND';
	}
	function isPrometheusError(error: Error) {
		return error.code === 'PROMETHEUS_NOT_FOUND';
	}

	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetchErrors(scope.uuid);
			mounted = true;
			refreshErrors(scope.uuid);
			iterateInformations();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !mounted}
	<PageLoading />
{:else}
	<main class="grid h-[calc(100vh_-_theme(spacing.16))] p-4">
		{#if $errorsStore.some((e) => isCephError(e) || isKubernetesError(e))}
			<div class="flex w-full flex-col items-center gap-4">
				{#if scope}
					{#if $errorsStore && $errorsStore.length > 0}
						{@const criticalErrors = $errorsStore.filter((e) => Number(e.level) > 3)}
						{#each criticalErrors as error}
							<Alert.Root variant="destructive">
								<Icon icon="material-symbols:warning-rounded" class="size-6" />
								<Alert.Title class="text-sm">{error.message}</Alert.Title>
								<Alert.Description class="text-xs text-destructive"
									>{error.details}</Alert.Description
								>
							</Alert.Root>
						{/each}
						{#if $currentInformationStore}
							<Alert.Root>
								<Icon icon="ph:info" class="size-6" />
								<Alert.Title class="flex items-center justify-between gap-2">
									<p class="text-sm">{$currentInformationStore.message}</p>
									<p class="text-sm">{$errorsStore.length} messages</p>
								</Alert.Title>
								<Alert.Description class="text-xs">
									{$currentInformationStore.details}
								</Alert.Description>
							</Alert.Root>
						{/if}
					{/if}
				{/if}

				<span class="flex h-full w-full items-center justify-evenly">
					{#each $errorsStore as error}
						{#if isCephError(error)}
							<CephError />
						{/if}
						{#if isKubernetesError(error)}
							<KubernetesError />
						{/if}
					{/each}
				</span>
			</div>
		{:else if $errorsStore.some((e) => isPrometheusError(e))}
			<PrometheusError />
		{:else}
			<Tabs.Root value="overview" class="w-full">
				<Tabs.List class="w-fit">
					<Tabs.Trigger value="overview">Overview</Tabs.Trigger>
					<Tabs.Trigger value="details">Details</Tabs.Trigger>
				</Tabs.List>

				<Tabs.Content value="overview">
					<div class="col-span-1">
						<div class="col-span-1 grid grid-cols-2 gap-2">
							<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
								<Card.Header>
									<Card.Title class="flex justify-between">
										<h1 class="text-3xl">Health</h1>
										<DropdownMenu.Root>
											<DropdownMenu.Trigger>
												<Button variant="ghost" class="flex items-start gap-2">
													<p>Alert</p>
													<Badge
														variant="destructive"
														class="flex h-4 w-4 justify-center rounded-full p-0 text-[12px] font-light"
													>
														{data.filter((h) => h.unhealth).length}
													</Badge>
												</Button>
											</DropdownMenu.Trigger>
											<DropdownMenu.Content>
												{#each data
													.filter((h) => h.unhealth)
													.sort((h) => h.unhealth) as healthDatum}
													<DropdownMenu.Item>
														<a
															href={healthDatum.link}
															class="flex w-full items-start justify-center gap-2"
														>
															<Badge
																variant="destructive"
																class="flex h-4 w-4 justify-center rounded-full p-0 text-[12px] font-light"
															>
																{healthDatum.unhealth}
															</Badge>
															{healthDatum.namespace}
														</a>
													</DropdownMenu.Item>
												{/each}
											</DropdownMenu.Content>
										</DropdownMenu.Root>
									</Card.Title>
									<Card.Description>
										<span class="text-sm font-light">Ratio of non-running pods over all pods</span>
									</Card.Description>
								</Card.Header>
								<Card.Content>
									<div class="flex h-full w-full items-center justify-center">
										<div
											class={cn(
												`h-[${Math.round(chartHeight.large * 1.732)}px] w-[${Math.round(chartHeight.large * 1.732)}px]`
											)}
										>
											<Chart>
												<Svg center>
													<Group y={chartHeight.large / 4}>
														<Arc
															value={(sum(data, (h) => h.unhealth) * 100) /
																sum(data, (h) => h.total)}
															domain={[0, 100]}
															outerRadius={chartHeight.large}
															innerRadius={-13}
															cornerRadius={13}
															range={[-120, 120]}
															class="fill-red-800"
															track={{ class: 'fill-red-50' }}
															let:value
														>
															<Text
																value={`${Math.round(value)}%`}
																textAnchor="middle"
																verticalAnchor="middle"
																class="text-5xl tabular-nums"
															/>
														</Arc>
													</Group>
												</Svg>
											</Chart>
										</div>
									</div>
								</Card.Content>
								<Card.Footer></Card.Footer>
							</Card.Root>
							<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
								<Card.Header>
									<Card.Title class="text-3xl">Storage</Card.Title>
									<Card.Description>
										<span class="text-sm font-light">Storage usage across all clusters</span>
									</Card.Description>
								</Card.Header>
								<Card.Content>
									<div class="flex h-full w-full items-center justify-center">
										<div
											class={cn(
												`h-[${Math.round(chartHeight.large * 1.732)}px] w-[${Math.round(chartHeight.large * 1.732)}px]`
											)}
										>
											<Chart>
												<Svg center>
													<Group y={chartHeight.large / 4}>
														<Arc
															value={Math.round((storages.used * 100) / storages.capacity)}
															domain={[0, 100]}
															outerRadius={chartHeight.large}
															innerRadius={-13}
															cornerRadius={13}
															range={[-120, 120]}
															class={metricColor((storages.used * 100) / storages.capacity)}
															track={{
																class: metricBackgroundColor(
																	(storages.used * 100) / storages.capacity
																)
															}}
															let:value
														>
															<Text
																value={`${value}%`}
																textAnchor="middle"
																verticalAnchor="middle"
																class="text-5xl tabular-nums"
															/>
															<Text
																value={`${storages.used} ${storages.unit} used over ${storages.capacity} ${storages.unit}`}
																textAnchor="middle"
																verticalAnchor="middle"
																dy={23 + 13}
																class="text-xs font-light"
															/>
														</Arc>
													</Group>
												</Svg>
											</Chart>
										</div>
									</div>
								</Card.Content>
							</Card.Root>
						</div>

						<div class="col-span-1 -mt-4 grid grid-cols-1 gap-2">
							<div class="col-span-1 grid grid-cols-3 gap-3">
								<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
									<Card.Header>
										<Card.Title class="text-2xl">CPU</Card.Title>
										<Card.Description
											><span class="text-sm font-light">CPU usage across all clusters</span
											></Card.Description
										>
									</Card.Header>
									<Card.Content>
										<div class="flex h-full w-full items-center justify-center">
											<div
												class={cn(
													`h-[${Math.round(chartHeight.small * 1.732)}px] w-[${Math.round(chartHeight.small * 1.732)}px]`
												)}
											>
												<Chart>
													<Svg center>
														<Group y={chartHeight.small / 4}>
															<Arc
																value={currentUsage.CPU}
																domain={[0, 100]}
																outerRadius={chartHeight.small}
																innerRadius={-13}
																cornerRadius={13}
																range={[-120, 120]}
																class={metricColor(currentUsage.CPU)}
																track={{
																	class: metricBackgroundColor(currentUsage.CPU)
																}}
															>
																<Text
																	value={`${currentUsage.CPU}%`}
																	textAnchor="middle"
																	verticalAnchor="middle"
																	class="text-4xl tabular-nums"
																/>
															</Arc>
														</Group>
													</Svg>
												</Chart>
											</div>
										</div>
									</Card.Content>
								</Card.Root>
								<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
									<Card.Header>
										<Card.Title class="text-2xl">GPU</Card.Title>
										<Card.Description
											><span class="text-sm font-light">GPU usage across all clusters</span
											></Card.Description
										>
									</Card.Header>
									<Card.Content>
										<div class="flex h-full w-full items-center justify-center">
											<div
												class={cn(
													`h-[${Math.round(chartHeight.small * 1.732)}px] w-[${Math.round(chartHeight.small * 1.732)}px]`
												)}
											>
												<Chart>
													<Svg center>
														<Group y={chartHeight.small / 4}>
															<Arc
																value={currentUsage.GPU}
																domain={[0, 100]}
																outerRadius={chartHeight.small}
																innerRadius={-13}
																cornerRadius={13}
																range={[-120, 120]}
																class={metricColor(currentUsage.GPU)}
																track={{
																	class: metricBackgroundColor(currentUsage.GPU)
																}}
															>
																<Text
																	value={`${currentUsage.GPU}%`}
																	textAnchor="middle"
																	verticalAnchor="middle"
																	class="text-4xl tabular-nums"
																/>
															</Arc>
														</Group>
													</Svg>
												</Chart>
											</div>
										</div>
									</Card.Content>
								</Card.Root>
								<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
									<Card.Header>
										<Card.Title class="text-2xl">Memory</Card.Title>
										<Card.Description
											><span class="text-sm font-light">Memory usage over all cluster</span
											></Card.Description
										>
									</Card.Header>
									<Card.Content>
										<div class="flex h-full w-full items-center justify-center">
											<div
												class={cn(
													`h-[${Math.round(chartHeight.small * 1.732)}px] w-[${Math.round(chartHeight.small * 1.732)}px]`
												)}
											>
												<Chart>
													<Svg center>
														<Group y={chartHeight.small / 4}>
															<Arc
																value={(currentUsage.memory.usage * 100) /
																	currentUsage.memory.capacity}
																domain={[0, 100]}
																outerRadius={chartHeight.small}
																innerRadius={-13}
																cornerRadius={13}
																range={[-120, 120]}
																class={metricColor(
																	(currentUsage.memory.usage * 100) / currentUsage.memory.capacity
																)}
																track={{
																	class: metricBackgroundColor(
																		(currentUsage.memory.usage * 100) / currentUsage.memory.capacity
																	)
																}}
															>
																<Text
																	value={`${Math.round((currentUsage.memory.usage * 100) / currentUsage.memory.capacity)}%`}
																	textAnchor="middle"
																	verticalAnchor="middle"
																	class="text-4xl tabular-nums"
																/>
															</Arc>
														</Group>
													</Svg>
												</Chart>
											</div>
										</div>
									</Card.Content>
								</Card.Root>
							</div>
						</div>
					</div>
				</Tabs.Content>

				<Tabs.Content value="details">
					<div class="col-span-1 grid grid-cols-2 gap-2">
						<div class="col-span-2">
							<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
								<Card.Header>
									<Card.Title class="text-3xl">Latency</Card.Title>
									<Card.Description>
										<span class="text-sm font-light">Latency trends for top 3 applications</span>
									</Card.Description>
								</Card.Header>
								<Card.Content>
									<div class="h-40 w-full resize overflow-visible">
										<AreaChart
											data={latencies}
											x="date"
											series={Object.keys(latencies[0])
												.filter((key) => key !== 'date')
												.map((k) => ({ key: k }))}
											legend
											props={{
												tooltip: {
													root: { class: 'bg-white/60 p-3 rounded shadow-lg' },
													header: { class: 'font-light' }
												}
											}}
											{renderContext}
											{debug}
										/>
									</div>
								</Card.Content>
							</Card.Root>
						</div>
						<div class="col-span-2 grid grid-cols-2 gap-2">
							<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
								<Card.Header>
									<Card.Title class="text-3xl">I/O</Card.Title>
									<Card.Description>
										<span class="text-sm font-light">Input/Output operations over time</span>
									</Card.Description>
								</Card.Header>
								<Card.Content>
									<div class="h-40 w-full resize overflow-visible">
										<AreaChart
											data={inputOutputs}
											x="date"
											series={[
												{ key: 'write', color: 'hsl(var(--color-info))' },
												{ key: 'read', color: 'hsl(var(--color-danger))' }
											]}
											legend
											props={{
												tooltip: {
													root: { class: 'bg-white/60 p-3 rounded shadow-lg' },
													header: { class: 'font-light' }
												}
											}}
											{renderContext}
											{debug}
										/>
									</div>
								</Card.Content>
							</Card.Root>
							<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
								<Card.Header>
									<Card.Title class="text-3xl">Resource</Card.Title>
									<Card.Description>
										<span class="text-sm font-light">Historical resource usage trends</span>
									</Card.Description>
								</Card.Header>
								<Card.Content>
									<div class="h-40 w-full resize overflow-visible">
										<AreaChart
											data={usages}
											x="date"
											series={[
												{ key: 'CPU', color: 'hsl(var(--color-info))' },
												{ key: 'GPU', color: 'hsl(var(--color-accent))' },
												{ key: 'Memory', color: 'hsl(var(--color-danger))' }
											]}
											legend
											props={{
												tooltip: {
													root: { class: 'bg-white/60 p-3 rounded shadow-lg' },
													header: { class: 'font-light' }
												}
											}}
											{renderContext}
											{debug}
										/>
									</div>
								</Card.Content>
							</Card.Root>
						</div>
					</div>
				</Tabs.Content>
			</Tabs.Root>
		{/if}
	</main>
{/if}
