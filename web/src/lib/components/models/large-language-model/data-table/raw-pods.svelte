<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row, Table as TableType } from '@tanstack/table-core';
	import { scaleUtc } from 'd3-scale';
	import { curveLinear } from 'd3-shape';
	import { LineChart } from 'layerchart';

	import type { Model } from '$lib/api/model/v1/model_pb';
	import * as Table from '$lib/components/custom/table/index.js';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Chart from '$lib/components/ui/chart';
	import { m } from '$lib/paraglide/messages';

	import { type LargeLanguageModel } from '../type';
	import type { Metrics } from '../types';
</script>

<script lang="ts">
	let {
		metrics,
		table,
		row
	}: { metrics: Metrics; table: TableType<Model>; row: Row<LargeLanguageModel> } = $props();

	const pods = $derived(row.original.pods ?? []);

	const total = $derived(pods.length);
	const running = $derived(pods.filter((p) => p.phase === 'Running').length);
	const ready = $derived(
		pods.filter((p) => (typeof p.ready === 'number' ? p.ready > 0 : Boolean(p.ready))).length
	);
	const averageRestarts = $derived(
		pods.map((p) => Number(p.restarts ?? 0)).reduce((p, n) => p + n, 0) / (pods.length || 1)
	);

	const phases = $derived(
		pods.reduce(
			(a, pod) => {
				a[pod.phase] = (a[pod.phase] ?? 0) + 1;
				return a;
			},
			{} as Record<string, number>
		)
	);
</script>

<Table.Row class="hover:[&,&>svelte-css-wrapper]:[&>th,td]:bg-transparent">
	<Table.Cell
		colspan={Object.keys(table.getHeaderGroups().flatMap((headerGroup) => headerGroup.headers))
			.length}
	>
		<div class="grid h-full grid-cols-3 items-start gap-8 p-4">
			<div class="col-span-1 flex h-full items-center justify-center p-4">
				<div class="grid w-fit grid-cols-3 gap-4">
					<div class="col-span-3 grid grid-cols-3 gap-12">
						<div class="col-span-1">
							<h3 class="text-xl font-bold">{total} Pods</h3>
						</div>
						<div class="col-span-2 flex items-center gap-1">
							<Icon icon="ph:heartbeat" class="size-6" />
							<div>
								<h6 class="text-xs text-muted-foreground">{m.phase()}</h6>
								<span class="flex flex-wrap gap-1">
									{#each Object.entries(phases) as [phaseName, count]}
										<Badge variant="outline" class="text-xs">
											{count}
											{phaseName}
										</Badge>
									{/each}
								</span>
							</div>
						</div>
					</div>
					<div class="col-span-3 grid grid-cols-3 gap-12">
						<div class="flex items-center gap-1">
							<Icon icon="ph:check-circle" class="size-6 text-chart-2" />
							<div>
								<h6 class="text-xs text-muted-foreground">{m.running()}</h6>
								<p class="text-base">{running}</p>
							</div>
						</div>
						<div class="flex items-center gap-1">
							<Icon icon="ph:stop-circle" class="size-6 text-chart-4" />
							<div>
								<h6 class="text-xs text-muted-foreground">{m.ready()}</h6>
								<p class="text-base">{ready}</p>
							</div>
						</div>
						<div class="flex items-center gap-1">
							<Icon icon="ph:warning-circle" class="size-6 text-chart-1" />
							<div>
								<h6 class="text-xs text-muted-foreground">{m.restarts()}</h6>
								<p class="text-base">{averageRestarts}</p>
							</div>
						</div>
					</div>
				</div>
			</div>
			<div class="col-span-2 rounded-lg border">
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>{m.pod()}</Table.Head>
							<Table.Head>{m.phase()}</Table.Head>
							<Table.Head>{m.ready()}</Table.Head>
							<Table.Head>{m.restarts()}</Table.Head>
							<Table.Head>{m.kv_cache()}</Table.Head>
							<Table.Head>{m.time_to_first_token()}</Table.Head>
							<Table.Head>{m.requests()}</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each row.original.pods as pod}
							<Table.Row>
								<Table.Cell>{pod.name}</Table.Cell>
								<Table.Cell>{pod.phase}</Table.Cell>
								<Table.Cell>{pod.ready}</Table.Cell>
								<Table.Cell>{pod.restarts}</Table.Cell>
								<Table.Cell>
									{@const configuration = {
										cache: { label: 'cache', color: 'var(--chart-2)' }
									} satisfies Chart.ChartConfig}
									<Chart.Container config={configuration} class="h-fit w-20">
										<LineChart
											data={pod.metrics.kv_cache}
											x="time"
											xScale={scaleUtc()}
											axis={false}
											series={[
												{
													key: 'value',
													label: configuration.cache.label,
													color: configuration.cache.color
												}
											]}
											props={{
												spline: {
													curve: curveLinear,
													motion: 'tween',
													strokeWidth: 2
												},
												xAxis: {
													format: (v: Date) =>
														v.toLocaleDateString('en-US', {
															month: 'short'
														})
												},
												highlight: { points: { r: 4 } }
											}}
										>
											{#snippet tooltip()}
												<Chart.Tooltip hideLabel>
													{#snippet formatter({ item, name, value })}
														<div
															style="--color-bg: {item.color}"
															class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
														></div>
														<div
															class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none"
														>
															<div class="grid gap-1.5">
																<span class="text-muted-foreground">{name}</span>
															</div>
															<p class="font-mono">
																{Number(value)}
															</p>
														</div>
													{/snippet}
												</Chart.Tooltip>
											{/snippet}
										</LineChart>
									</Chart.Container>
								</Table.Cell>
								<Table.Cell>
									{@const configuration = {
										time: { label: 'time', color: 'var(--chart-1)' }
									} satisfies Chart.ChartConfig}
									<Chart.Container config={configuration} class="h-fit w-20">
										<LineChart
											data={pod.metrics.kv_cache}
											x="time"
											xScale={scaleUtc()}
											axis={false}
											series={[
												{
													key: 'value',
													label: configuration.time.label,
													color: configuration.time.color
												}
											]}
											props={{
												spline: {
													curve: curveLinear,
													motion: 'tween',
													strokeWidth: 2
												},
												xAxis: {
													format: (v: Date) =>
														v.toLocaleDateString('en-US', {
															month: 'short'
														})
												},
												highlight: { points: { r: 4 } }
											}}
										>
											{#snippet tooltip()}
												<Chart.Tooltip hideLabel>
													{#snippet formatter({ item, name, value })}
														<div
															style="--color-bg: {item.color}"
															class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
														></div>
														<div
															class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none"
														>
															<div class="grid gap-1.5">
																<span class="text-muted-foreground">{name}</span>
															</div>
															<p class="font-mono">
																{Number(value)}
															</p>
														</div>
													{/snippet}
												</Chart.Tooltip>
											{/snippet}
										</LineChart>
									</Chart.Container>
								</Table.Cell>
								<Table.Cell>
									{@const configuration = {
										request: { label: 'request', color: 'var(--chart-1)' }
									} satisfies Chart.ChartConfig}
									<Chart.Container config={configuration} class="h-fit w-20">
										<LineChart
											data={pod.metrics.kv_cache}
											x="time"
											xScale={scaleUtc()}
											axis={false}
											series={[
												{
													key: 'value',
													label: configuration.request.label,
													color: configuration.request.color
												}
											]}
											props={{
												spline: {
													curve: curveLinear,
													motion: 'tween',
													strokeWidth: 2
												},
												xAxis: {
													format: (v: Date) =>
														v.toLocaleDateString('en-US', {
															month: 'short'
														})
												},
												highlight: { points: { r: 4 } }
											}}
										>
											{#snippet tooltip()}
												<Chart.Tooltip hideLabel>
													{#snippet formatter({ item, name, value })}
														<div
															style="--color-bg: {item.color}"
															class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
														></div>
														<div
															class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none"
														>
															<div class="grid gap-1.5">
																<span class="text-muted-foreground">{name}</span>
															</div>
															<p class="font-mono">
																{Number(value)}
															</p>
														</div>
													{/snippet}
												</Chart.Tooltip>
											{/snippet}
										</LineChart>
									</Chart.Container>
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
		</div>
	</Table.Cell>
</Table.Row>
