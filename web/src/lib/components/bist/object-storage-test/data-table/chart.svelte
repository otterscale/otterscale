<script lang="ts" generics="TData">
	import type { TestResult } from '$gen/api/bist/v1/bist_pb';
	import { BistDashboardManager } from '$lib/components/bist/utils/bistManager';
	import { Chart as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Template from '$lib/components/dashboard/utils/templates';
	import * as Select from "$lib/components/ui/select/index.js";
	import { formatCapacityV2 as formatCapacity } from '$lib/formatter';
	import { type Table } from '@tanstack/table-core';
	import { scaleLog } from 'd3-scale';
	import dayjs from 'dayjs';
	import { ScatterChart, Tooltip } from 'layerchart';

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

    // 確保泛型 TData 被正確設定為 TestResult
    let { table }: { table: Table<TestResult> } = $props();

	const dashboardManager = new BistDashboardManager(table);


    const modes = [
        { value: "get", label: "GET" },
        { value: "put", label: "PUT" },
        { value: "delete", label: "DELETE" }
    ];
    
    let modeThroughputFastest = $state("get");
    let modeThroughputMedian = $state("get");

    const triggerBandwidthContent = $derived(
        modes.find((m) => m.value === modeThroughputFastest)?.label ?? "Select a mode"
    );
    const triggerIOContent = $derived(
        modes.find((m) => m.value === modeThroughputMedian)?.label ?? "Select a mode"
    );
</script>

<Layout>
    {@const { get: getTmp, put: putTmp, delete: deleteTmp } = dashboardManager.getWarpOutputs()}
    <!-- Show FioOutputs at the side -->
	<Template.Area title="Throughput">
		{#snippet hint()}
			<p>Throughput Fastest</p>
		{/snippet}
        {#snippet controller()}
            <Select.Root type="single" name="ioMode" bind:value={modeThroughputFastest}>
            <Select.Trigger class="w-[180px]">
                {triggerBandwidthContent}
            </Select.Trigger>
            <Select.Content>
                <Select.Group>
                {#each modes as mode (mode.value)}
                    <Select.Item
                    value={mode.value}
                    label={mode.label}
                    >
                    {mode.label}
                    </Select.Item>
                {/each}
                </Select.Group>
            </Select.Content>
            </Select.Root>
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
                            modeThroughputFastest === "get" ? getTmp :
                            modeThroughputFastest === "put" ? putTmp :
                            deleteTmp
                        )}
						yScale={scaleLog()}
						props={{
							xAxis: { 
								tweened: { duration: 200 },
								format: (d: Date) => dayjs(d).format('MM/DD')
							},
							yAxis: {
								format: (v: number) => {
									const capacity = formatCapacity(v);
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
									value={`${Number(formatCapacity(data.bytesFastest).value).toFixed(0)} ${formatCapacity(data.bytesFastest).unit}`}
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
        {#snippet controller()}
            <Select.Root type="single" name="ioMode" bind:value={modeThroughputMedian}>
            <Select.Trigger class="w-[180px]">
                {triggerIOContent}
            </Select.Trigger>
            <Select.Content>
                <Select.Group>
                {#each modes as mode (mode.value)}
                    <Select.Item
                    value={mode.value}
                    label={mode.label}
                    >
                    {mode.label}
                    </Select.Item>
                {/each}
                </Select.Group>
            </Select.Content>
            </Select.Root>
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
                            modeThroughputMedian === "get" ? getTmp :
                            modeThroughputMedian === "put" ? putTmp :
                            deleteTmp
                        )}
						yScale={scaleLog()}
						props={{
							xAxis: { 
								tweened: { duration: 200 },
								format: (d: Date) => dayjs(d).format('MM/DD')
							},
							yAxis: {
								format: (v: number) => {
									const capacity = formatCapacity(v);
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
									value={`${Number(formatCapacity(data.bytesMedian).value).toFixed(0)} ${formatCapacity(data.bytesMedian).unit}`}
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