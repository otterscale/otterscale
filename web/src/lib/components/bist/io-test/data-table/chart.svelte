<script lang="ts" generics="TData">
	import type { TestResult } from '$gen/api/bist/v1/bist_pb';
	import { BistDashboardManager } from '$lib/components/bist/utils/bistManager';
	import { Chart as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Template from '$lib/components/dashboard/utils/templates';
	import { formatByte, formatLatencyNano } from '$lib/formatter';
	import { type Table } from '@tanstack/table-core';
	import { scaleLog } from 'd3-scale';
	import dayjs from 'dayjs';
	import { ScatterChart, Tooltip } from 'layerchart';
	import Pickers from './pickers.svelte';

	let { table }: { table: Table<TestResult> } = $props();
	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;
	let mode = $state('read');
	const dashboardManager = new BistDashboardManager(table);
</script>

<Pickers bind:selectedMode={mode} />
<Layout>
	{@const { read: readTmp, write: writeTmp, trim: trimTmp } = dashboardManager.getFioOutputs()}
	<!-- Show FioOutputs at the side -->
	{console.log('readTmp', readTmp)}
	<Template.Area title="Bandwidth">
		{#snippet hint()}
			<p>Bandwidth Bytes</p>
		{/snippet}
		{#snippet content()}
			<div class="h-[200px] w-full resize overflow-visible">
				<ScatterChart
					x="completedAt"
					y="bandwidthBytes"
					series={Object.values(mode === 'read' ? readTmp : mode === 'write' ? writeTmp : trimTmp)}
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
						points: { tweened: { duration: 200 } }
					}}
					legend={{
						classes: { root: '-mb-[50px] w-full overflow-auto' }
					}}
					{renderContext}
					{debug}
				>
					<svelte:fragment slot="tooltip">
						<Tooltip.Root let:data>
							<Tooltip.Header class="font-light">{data.seriesKey}</Tooltip.Header>
							<Tooltip.List>
								<Tooltip.Item label="Name" value={data.name} />
								<Tooltip.Item
									label="Bandwidth"
									value={`${Number(formatByte(data.bandwidthBytes).value).toFixed(0)} ${formatByte(data.bandwidthBytes).unit}`}
								/>
								<Tooltip.Item
									label="Date"
									value={dayjs(data.completedAt).format('YYYY/MM/DD HH:mm')}
								/>
							</Tooltip.List>
						</Tooltip.Root>
					</svelte:fragment>
				</ScatterChart>
			</div>
		{/snippet}
	</Template.Area>

	<Template.Area title="IOPS">
		{#snippet hint()}
			<p>IO Per Second</p>
		{/snippet}
		{#snippet content()}
			<div class="h-[200px] w-full resize overflow-visible">
				<ScatterChart
					x="completedAt"
					y="ioPerSecond"
					series={Object.values(mode === 'read' ? readTmp : mode === 'write' ? writeTmp : trimTmp)}
					yScale={scaleLog()}
					props={{
						xAxis: {
							tweened: { duration: 200 },
							format: (d: Date) => dayjs(d).format('MM/DD')
						},
						grid: { tweened: { duration: 200 } },
						points: { tweened: { duration: 200 } }
					}}
					legend={{
						classes: { root: '-mb-[50px] w-full overflow-auto' }
					}}
					{renderContext}
					{debug}
				>
					<svelte:fragment slot="tooltip">
						<Tooltip.Root let:data>
							<Tooltip.Header class="font-light">{data.seriesKey}</Tooltip.Header>
							<Tooltip.List>
								<Tooltip.Item label="Name" value={data.name} />
								<Tooltip.Item label="IO" value={data.ioPerSecond} />
								<Tooltip.Item
									label="Date"
									value={dayjs(data.completedAt).format('YYYY/MM/DD HH:mm')}
								/>
							</Tooltip.List>
						</Tooltip.Root>
					</svelte:fragment>
				</ScatterChart>
			</div>
		{/snippet}
	</Template.Area>

	<Template.Area title="Latency">
		{#snippet hint()}
			<p>Mean Latency</p>
		{/snippet}
		{#snippet content()}
			<div class="h-[200px] w-full resize overflow-visible">
				<ScatterChart
					x="completedAt"
					y="latency"
					series={Object.values(mode === 'read' ? readTmp : mode === 'write' ? writeTmp : trimTmp)}
					props={{
						xAxis: {
							tweened: { duration: 200 },
							format: (d: Date) => dayjs(d).format('MM/DD')
						},
						yAxis: {
							format: (v: number) => {
								const latency = formatLatencyNano(v);
								return `${Number(latency.value).toFixed(0)} ${latency.unit}`;
							}
						},
						grid: { tweened: { duration: 200 } },
						points: { tweened: { duration: 200 } }
					}}
					legend={{
						classes: { root: '-mb-[50px] w-full overflow-auto' }
					}}
					{renderContext}
					{debug}
				>
					<svelte:fragment slot="tooltip">
						<Tooltip.Root let:data>
							<Tooltip.Header class="font-light">{data.seriesKey}</Tooltip.Header>
							<Tooltip.List>
								<Tooltip.Item label="Name" value={data.name} />
								<Tooltip.Item
									label="Latency"
									value={`${Number(formatLatencyNano(data.latency).value).toFixed(0)} ${formatLatencyNano(data.latency).unit}`}
								/>
								<Tooltip.Item
									label="Date"
									value={dayjs(data.completedAt).format('YYYY/MM/DD HH:mm')}
								/>
							</Tooltip.List>
						</Tooltip.Root>
					</svelte:fragment>
				</ScatterChart>
			</div>
		{/snippet}
	</Template.Area>
</Layout>
