<script lang="ts">
	const healthData = group(healthRawData, (d) => d.type).get('error') ?? [];

	import { cn } from '$lib/utils';
	import { type Error, type Machine } from '$gen/api/nexus/v1/nexus_pb';
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

	import PrometheusError from './error/prometheus.svelte';
	import KubernetesError from './error/kubernetes.svelte';
	import CephError from './error/ceph.svelte';

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	let { errors }: { errors: Error[] } = $props();

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

	let data = $state(healthData);

	function isCephError(error: Error) {
		return error.code === 'CEPH_NOT_FOUND';
	}
	function isKubernetesError(error: Error) {
		return error.code === 'KUBERNETES_NOT_FOUND';
	}
	function isPrometheusError(error: Error) {
		return error.code === 'PROMETHEUS_NOT_FOUND';
	}
</script>

<main class="flex h-[calc(100vh_-_theme(spacing.16))] justify-between gap-4 border p-4">
	{#if errors.some((e) => isCephError(e) || isKubernetesError(e))}
		<div class="flex w-full justify-center">
			<span class="flex justify-between">
				{#each errors as error}
					{#if isCephError(error)}
						<CephError />
					{/if}
					{#if isKubernetesError(error)}
						<KubernetesError />
					{/if}
				{/each}
			</span>
		</div>
	{:else if errors.some((e) => isPrometheusError(e))}
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
											{#each data.filter((h) => h.unhealth).sort((h) => h.unhealth) as healthDatum}
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
											`h-[calc(${chartHeight.large}px*1.732)] w-[calc(${chartHeight.large}px*1.732)]`
										)}
									>
										<Chart>
											<Svg center>
												<Group y={chartHeight.large / 4}>
													<Arc
														value={(sum(data, (h) => h.unhealth) * 100) / sum(data, (h) => h.total)}
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
											`h-[calc(${chartHeight.large}px*1.732)] w-[calc(${chartHeight.large}px*1.732)]`
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
												`h-[calc(${chartHeight.small}px*1.732)] w-[calc(${chartHeight.small}px*1.732)]`
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
												`h-[calc(${chartHeight.small}px*1.732)] w-[calc(${chartHeight.small}px*1.732)]`
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
												`h-[calc(${chartHeight.small}px*1.732)] w-[calc(${chartHeight.small}px*1.732)]`
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
