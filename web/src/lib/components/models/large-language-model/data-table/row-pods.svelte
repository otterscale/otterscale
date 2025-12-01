<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row, Table as TableType } from '@tanstack/table-core';
	import { scaleUtc } from 'd3-scale';
	import { LineChart } from 'layerchart';
	import { SampleValue } from 'prometheus-query';

	import type { Model } from '$lib/api/model/v1/model_pb';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import * as Chart from '$lib/components/ui/chart';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { formatTimeAgo } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	import type { Metrics } from '../types';
	import Log from './action-log.svelte';
</script>

<script lang="ts">
	let {
		metrics,
		table,
		row,
		scope,
		namespace
	}: {
		metrics: Metrics;
		table: TableType<Model>;
		row: Row<Model>;
		scope: string;
		namespace: string;
	} = $props();

	const pods = $derived(row.original.pods ?? []);

	const total = $derived(pods.length);
	const running = $derived(pods.filter((pod) => pod.phase === 'Running').length);
	const averageRestarts = $derived(
		total > 0 ? pods.reduce((a, pod) => a + Number(pod.restarts), 0) / total : 0
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
	const failed = $derived(pods.filter((pod) => pod?.lastCondition?.status !== 'True').length);
</script>

<Table.Row class="bg-muted/50 hover:[&,&>svelte-css-wrapper]:[&>th,td]:bg-transparent">
	<Table.Cell
		colspan={Object.keys(table.getHeaderGroups().flatMap((headerGroup) => headerGroup.headers))
			.length}
	>
		<div class="h-full space-y-4 border-l-4 border-border px-8 pt-8 pb-12">
			<div class="flex items-center justify-between gap-4 p-4">
				<h3 class="text-2xl font-bold">{total} Pods</h3>
				<div class="flex items-center gap-4">
					<div class="flex items-center gap-1">
						<Icon icon="ph:heartbeat-duotone" class="size-6" />
						<div>
							<h6 class="text-xs text-muted-foreground">{m.phase()}</h6>
							<span class="flex flex-wrap gap-1">
								{#each Object.entries(phases) as [phaseName, count], index (index)}
									<Badge variant="outline" class="text-xs">
										{count}
										{phaseName}
									</Badge>
								{/each}
							</span>
						</div>
					</div>
					<div class="flex items-center gap-1">
						<Icon icon="ph:check-circle-duotone" class="size-6 text-chart-2" />
						<div>
							<h6 class="text-xs text-muted-foreground">{m.running()}</h6>
							<p class="text-base">{running}</p>
						</div>
					</div>
					<div class="flex items-center gap-1">
						<Icon icon="ph:stop-circle-duotone" class="size-6 text-chart-4" />
						<div>
							<h6 class="text-xs text-muted-foreground">{m.restarts()}</h6>
							<p class="text-base">{averageRestarts}</p>
						</div>
					</div>
					<div class="flex items-center gap-1">
						<Icon icon="ph:pause-circle-duotone" class="size-6 text-chart-1" />
						<div>
							<h6 class="text-xs text-muted-foreground">{m.failed()}</h6>
							<p class="text-base">{failed}</p>
						</div>
					</div>
				</div>
			</div>
			<div class="rounded-lg border">
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>{m.pod()}</Table.Head>
							<Table.Head>{m.phase()}</Table.Head>
							<Table.Head class="text-end">{m.ready()}</Table.Head>
							<Table.Head class="text-end">{m.restarts()}</Table.Head>
							<Table.Head class="text-start">{m.last_condition()}</Table.Head>
							<Table.Head class="text-center">{m.kv_cache()}</Table.Head>
							<Table.Head class="text-center">{m.time_to_first_token()}</Table.Head>
							<Table.Head class="text-center">{m.requests()}</Table.Head>
							<Table.Head class="text-end">{m.log()}</Table.Head>
							<Table.Head class="text-end">{m.create_time()}</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each row.original.pods as pod (pod.name)}
							<Table.Row>
								<Table.Cell>{pod.name}</Table.Cell>
								<Table.Cell>
									<Badge variant="outline">
										{pod.phase}
									</Badge>
								</Table.Cell>
								<Table.Cell class="text-end">
									{pod.ready}
								</Table.Cell>
								<Table.Cell class="text-end">{pod.restarts}</Table.Cell>
								<Table.Cell class="text-start">
									{#if pod.lastCondition}
										{@const status = pod.lastCondition.status}
										{#if status === 'True'}
											<Badge variant="outline">{pod.lastCondition.type}</Badge>
										{:else}
											<div class="flex items-center gap-1 text-destructive">
												<Badge variant="destructive">{pod.lastCondition.reason}</Badge>
												<Tooltip.Provider>
													<Tooltip.Root>
														<Tooltip.Trigger>
															<p class="max-w-[100px] truncate">
																{pod.lastCondition.message}
															</p>
														</Tooltip.Trigger>
														<Tooltip.Content>
															{pod.lastCondition.message}
														</Tooltip.Content>
													</Tooltip.Root>
												</Tooltip.Provider>
											</div>
										{/if}
									{/if}
								</Table.Cell>
								<Table.Cell class="text-center">
									{@const configuration = {
										cache: { label: 'cache', color: 'var(--chart-2)' }
									} satisfies Chart.ChartConfig}
									{@const kvCaches: SampleValue[] = metrics.kvCache.get(pod.name) ?? []}
									{#if kvCaches.length > 0}
										<Chart.Container config={configuration} class="h-10 w-full">
											<LineChart
												data={kvCaches}
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
														strokeWidth: 2
													}
												}}
												grid={false}
											>
												{#snippet tooltip()}
													<Chart.Tooltip hideLabel>
														{#snippet formatter({ item, name, value })}
															<div
																class="flex flex-1 shrink-0 items-center justify-start gap-1 font-mono text-xs leading-none"
																style="--color-bg: {item.color}"
															>
																<Icon icon="ph:square-fill" class="text-(--color-bg)" />
																<h1 class="font-semibold text-muted-foreground">{name}</h1>
																<p class="ml-auto">{(Number(value) * 100).toFixed(2)} %</p>
															</div>
														{/snippet}
													</Chart.Tooltip>
												{/snippet}
											</LineChart>
										</Chart.Container>
									{/if}
								</Table.Cell>
								<Table.Cell class="text-center">
									{@const configuration = {
										time: { label: 'time', color: 'var(--chart-1)' }
									} satisfies Chart.ChartConfig}
									{@const timeToFirstTokens: SampleValue[] = metrics.timeToFirstToken.get(pod.name) ?? []}
									{#if timeToFirstTokens.length > 0}
										<Chart.Container config={configuration} class="h-10 w-full">
											<LineChart
												data={timeToFirstTokens}
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
														strokeWidth: 2
													}
												}}
												grid={false}
											>
												{#snippet tooltip()}
													<Chart.Tooltip hideLabel>
														{#snippet formatter({ item, name, value })}
															<div
																class="flex flex-1 shrink-0 items-center justify-start gap-1 font-mono text-xs leading-none"
																style="--color-bg: {item.color}"
															>
																<Icon icon="ph:square-fill" class="text-(--color-bg)" />
																<h1 class="font-semibold text-muted-foreground">{name}</h1>
																<p class="ml-auto">{Number(value).toFixed(2)} sec.</p>
															</div>
														{/snippet}
													</Chart.Tooltip>
												{/snippet}
											</LineChart>
										</Chart.Container>
									{/if}
								</Table.Cell>
								<Table.Cell class="text-center">
									{@const configuration = {
										time: { label: 'time', color: 'var(--chart-1)' }
									} satisfies Chart.ChartConfig}
									{@const requestLatencies: SampleValue[] = metrics.requestLatency.get(pod.name) ?? []}
									<Chart.Container config={configuration} class="h-10 w-full">
										<LineChart
											data={requestLatencies}
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
													strokeWidth: 2
												}
											}}
											grid={false}
										>
											{#snippet tooltip()}
												<Chart.Tooltip hideLabel>
													{#snippet formatter({ item, name, value })}
														<div
															class="flex flex-1 shrink-0 items-center justify-start gap-1 font-mono text-xs leading-none"
															style="--color-bg: {item.color}"
														>
															<Icon icon="ph:square-fill" class="text-(--color-bg)" />
															<h1 class="font-semibold text-muted-foreground">{name}</h1>
															<p class="ml-auto">{Number(value).toFixed(2)} sec.</p>
														</div>
													{/snippet}
												</Chart.Tooltip>
											{/snippet}
										</LineChart>
									</Chart.Container>
								</Table.Cell>
								<Table.Cell class="text-end">
									{#if pod}
										<Log {pod} {scope} {namespace} />
									{/if}
								</Table.Cell>
								<Table.Cell class="text-end">
									{#if pod.createdAt}
										<Tooltip.Provider>
											<Tooltip.Root>
												<Tooltip.Trigger>
													{formatTimeAgo(timestampDate(pod.createdAt))}
												</Tooltip.Trigger>
												<Tooltip.Content>
													{timestampDate(pod.createdAt)}
												</Tooltip.Content>
											</Tooltip.Root>
										</Tooltip.Provider>
									{/if}
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
		</div>
	</Table.Cell>
</Table.Row>
