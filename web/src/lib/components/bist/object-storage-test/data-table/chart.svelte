<script lang="ts" generics="TData">
	import type { TestResult } from '$gen/api/bist/v1/bist_pb';
	import { BistDashboardManager } from '$lib/components/bist/utils/bistManager';
	import { Chart as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Template from '$lib/components/dashboard/utils/templates';
	import { formatByte } from '$lib/formatter';
	import { type Table } from '@tanstack/table-core';
	import dayjs from 'dayjs';
	import { ScatterChart, Tooltip } from 'layerchart';
	import Pickers from './pickers.svelte';

    let { table }: { table: Table<TestResult> } = $props();
	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;
	let mode = $state("put");
	const dashboardManager = new BistDashboardManager(table);
</script>

<Pickers bind:selectedMode={mode} />
<Layout>
    {@const { get: getTmp, put: putTmp, delete: deleteTmp } = dashboardManager.getWarpOutputs()}
    <!-- Show FioOutputs at the side -->
	<Template.Area title="Throughput">
		{#snippet hint()}
			<p>Throughput Fastest</p>
		{/snippet}
		{#snippet description()}
			<p class="text-xl">Fastest</p>
		{/snippet}
		{#snippet content()}
				<div class="h-[200px] w-full resize overflow-visible">
					<ScatterChart
						x="completedAt"
						y="bytesFastest"
                        series={Object.values(
                            mode === "get" ? getTmp :
                            mode === "put" ? putTmp :
                            deleteTmp
                        )}
						props={{
							xAxis: { 
								tweened: { duration: 200 },
								format: (d: Date) => dayjs(d).format('MM/DD')
							},
							yAxis: {
								format: (v: number) => {
									const capacity = formatByte(v);
									return `${Number(capacity.value).toFixed(0)} ${capacity.unit}`;
								}
							},
							grid: { tweened: { duration: 200 } },
							points: { tweened: { duration: 200 } },
						}}
						legend={{
							classes: { root: '-mb-[50px] w-full overflow-auto' }
						}}
						{renderContext}
						{debug}
					>
						<svelte:fragment slot="tooltip">
							<Tooltip.Root let:data>
							<Tooltip.Header class='font-light'>{data.seriesKey}</Tooltip.Header>
							<Tooltip.List>
								<Tooltip.Item label="Name" value={(data.name)} />
								<Tooltip.Item
									label="Bandwidth"
									value={`${Number(formatByte(data.bytesFastest).value).toFixed(0)} ${formatByte(data.bytesFastest).unit}`}
								/>
								<Tooltip.Item label="Date" value={dayjs(data.completedAt).format('YYYY/MM/DD HH:mm')} />
							</Tooltip.List>
							</Tooltip.Root>
						</svelte:fragment>
					</ScatterChart>
				</div>
		{/snippet}
	</Template.Area>

	<Template.Area title="Throughput">
		{#snippet hint()}
			<p>Throughput Slowest</p>
		{/snippet}
		{#snippet description()}
			<p class="text-xl">Slowest</p>
		{/snippet}
		{#snippet content()}
				<div class="h-[200px] w-full resize overflow-visible">
					<ScatterChart
						x="completedAt"
						y="bytesSlowest"
                        series={Object.values(
                            mode === "get" ? getTmp :
                            mode === "put" ? putTmp :
                            deleteTmp
                        )}
						props={{
							xAxis: { 
								tweened: { duration: 200 },
								format: (d: Date) => dayjs(d).format('MM/DD')
							},
							yAxis: {
								format: (v: number) => {
									const capacity = formatByte(v);
									return `${Number(capacity.value).toFixed(0)} ${capacity.unit}`;
								}
							},
							grid: { tweened: { duration: 200 } },
							points: { tweened: { duration: 200 } },
						}}
						legend={{
							classes: { root: '-mb-[50px] w-full overflow-auto' }
						}}
						{renderContext}
						{debug}
					>
						<svelte:fragment slot="tooltip">
							<Tooltip.Root let:data>
							<Tooltip.Header class='font-light'>{data.seriesKey}</Tooltip.Header>
							<Tooltip.List>
								<Tooltip.Item label="Name" value={(data.name)} />
								<Tooltip.Item
									label="Bandwidth"
									value={`${Number(formatByte(data.bytesSlowest).value).toFixed(0)} ${formatByte(data.bytesSlowest).unit}`}
								/>
								<Tooltip.Item label="Date" value={dayjs(data.completedAt).format('YYYY/MM/DD HH:mm')} />
							</Tooltip.List>
							</Tooltip.Root>
						</svelte:fragment>
					</ScatterChart>
				</div>
		{/snippet}
	</Template.Area>

	<Template.Area title="Throughput">
		{#snippet hint()}
			<p>Throughput Fastest</p>
		{/snippet}
		{#snippet description()}
			<p class="text-xl">Median</p>
		{/snippet}
		{#snippet content()}
				<div class="h-[200px] w-full resize overflow-visible">
					<ScatterChart
						x="completedAt"
						y="bytesMedian"
                        series={Object.values(
                            mode === "get" ? getTmp :
                            mode === "put" ? putTmp :
                            deleteTmp
                        )}
						props={{
							xAxis: { 
								tweened: { duration: 200 },
								format: (d: Date) => dayjs(d).format('MM/DD')
							},
							yAxis: {
								format: (v: number) => {
									const capacity = formatByte(v);
									return `${Number(capacity.value).toFixed(0)} ${capacity.unit}`;
								}
							},
							grid: { tweened: { duration: 200 } },
							points: { tweened: { duration: 200 } },
						}}
						legend={{
							classes: { root: '-mb-[50px] w-full overflow-auto' }
						}}
						{renderContext}
						{debug}
					>
						<svelte:fragment slot="tooltip">
							<Tooltip.Root let:data>
							<Tooltip.Header class='font-light'>{data.seriesKey}</Tooltip.Header>
							<Tooltip.List>
								<Tooltip.Item label="Name" value={(data.name)} />
								<Tooltip.Item
									label="Bandwidth"
									value={`${Number(formatByte(data.bytesMedian).value).toFixed(0)} ${formatByte(data.bytesMedian).unit}`}
								/>
								<Tooltip.Item label="Date" value={dayjs(data.completedAt).format('YYYY/MM/DD HH:mm')} />
							</Tooltip.List>
							</Tooltip.Root>
						</svelte:fragment>
					</ScatterChart>
				</div>
		{/snippet}
	</Template.Area>
</Layout>