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
	const succeeded = $derived(pods.filter((pod) => pod.phase === 'Succeeded').length);
	const failed = $derived(pods.filter((pod) => pod.phase === 'Failed').length);
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
		<div class="m-4 h-full space-y-4 border-l-4 border-border p-8">
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
						<Icon icon="ph:play-circle-duotone" class="size-6 text-chart-4" />
						<div>
							<h6 class="text-xs text-muted-foreground">{m.running()}</h6>
							<p class="text-base">{running}</p>
						</div>
					</div>
					<div class="flex items-center gap-1">
						<Icon icon="ph:check-circle-duotone" class="size-6 text-chart-2" />
						<div>
							<h6 class="text-xs text-muted-foreground">{m.succeeded()}</h6>
							<p class="text-base">{succeeded}</p>
						</div>
					</div>
					<div class="flex items-center gap-1">
						<Icon icon="ph:x-circle-duotone" class="size-6 text-chart-1" />
						<div>
							<h6 class="text-xs text-muted-foreground">{m.failed()}</h6>
							<p class="text-base">{failed}</p>
						</div>
					</div>
				</div>
			</div>
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>{m.pod()}</Table.Head>
						<Table.Head>{m.phase()}</Table.Head>
						<Table.Head class="text-end">{m.ready()}</Table.Head>
						<Table.Head class="text-end">{m.restarts()}</Table.Head>
						<Table.Head class="text-start">{m.conditions()}</Table.Head>
						<Table.Head class="text-center">{m.time_to_first_token()}</Table.Head>
						<Table.Head class="text-center">{m.request_latency()}</Table.Head>
						<Table.Head class="text-center">{m.log()}</Table.Head>
						<Table.Head class="text-end">{m.create_time()}</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#if row.original.pods.length > 0}
						{#each row.original.pods as pod (pod.name)}
							<Table.Row>
								<Table.Cell>{pod.name}</Table.Cell>
								<Table.Cell>
									{pod.phase}
								</Table.Cell>
								<Table.Cell class="text-end">
									{pod.ready}
								</Table.Cell>
								<Table.Cell class="text-end">{pod.restarts}</Table.Cell>
								<Table.Cell class="text-start">
									{#if pod.conditions}
										{@const trueConditions = pod.conditions.filter(
											(condition) => condition.status === 'True'
										)}
										<div class="flex flex-wrap gap-1">
											{#each trueConditions as trueCondition, index (index)}
												<Tooltip.Provider>
													<Tooltip.Root>
														<Tooltip.Trigger>
															<Badge variant="outline">{trueCondition.type}</Badge>
														</Tooltip.Trigger>
														<Tooltip.Content>
															{#if trueCondition.message}
																{trueCondition.message}
															{:else}
																{trueCondition.type}
															{/if}
														</Tooltip.Content>
													</Tooltip.Root>
												</Tooltip.Provider>
											{/each}
										</div>
									{/if}
								</Table.Cell>
								<Table.Cell class="text-center">
									{@const configuration = {
										time: { label: 'time', color: 'var(--chart-1)' }
									} satisfies Chart.ChartConfig}
									{@const timeToFirstTokens: SampleValue[] = metrics.timeToFirstToken?.get(pod.name) ?? []}
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
										time: { label: 'latency', color: 'var(--chart-1)' }
									} satisfies Chart.ChartConfig}
									{@const requestLatencies = (metrics.requestLatency?.get(pod.name) ?? []).map(
										(sampleValue: SampleValue) => ({
											time: sampleValue.time,
											value: sampleValue.value && !isNaN(sampleValue.value) ? sampleValue.value : 0
										})
									)}
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
								<Table.Cell>
									{#if pod}
										<div class="flex justify-center">
											<Log {pod} {scope} {namespace} />
										</div>
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
					{:else}
						<Table.Row>
							<Table.Cell colspan={10}>
								<Table.Empty />
							</Table.Cell>
						</Table.Row>
					{/if}
				</Table.Body>
			</Table.Root>
		</div>
	</Table.Cell>
</Table.Row>
