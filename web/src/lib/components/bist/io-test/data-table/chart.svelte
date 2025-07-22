<script lang="ts" generics="TData">
	import type { TestResult } from '$gen/api/bist/v1/bist_pb';
	import { BistDashboardManager } from '$lib/components/bist/utils/bistManager';
	import { type Table } from '@tanstack/table-core';
	
	import { Chart as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Template from '$lib/components/dashboard/utils/templates';
	import { formatCapacityV2 as formatCapacity } from '$lib/formatter';
	import { scaleLog } from 'd3-scale';
	import dayjs from 'dayjs';
	import { ScatterChart, Tooltip } from 'layerchart';
    import * as Select from "$lib/components/ui/select/index.js";

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;


    // 確保泛型 TData 被正確設定為 TestResult
    let { table }: { table: Table<TestResult> } = $props();

	const dashboardManager = new BistDashboardManager(table);


    const modes = [
        { value: "read", label: "Read" },
        { value: "write", label: "Write" },
        { value: "trim", label: "Trim" }
    ];
    
    let modeBandwidth = $state("read");
    let modeIO = $state("read");
    let modeLatency = $state("read");

    const triggerBandwidthContent = $derived(
        modes.find((m) => m.value === modeBandwidth)?.label ?? "Select a mode"
    );
    const triggerIOContent = $derived(
        modes.find((m) => m.value === modeIO)?.label ?? "Select a mode"
    );
    const triggerLatencyContent = $derived(
        modes.find((m) => m.value === modeLatency)?.label ?? "Select a mode"
    );
</script>

<Layout>
    {@const { read: readTmp, write: writeTmp, trim: trimTmp } = dashboardManager.getFioOutputs()}
    <!-- Show FioOutputs at the side -->
	<Template.Area title="Bandwidth">
		{#snippet hint()}
			<p>Bandwidth Bytes</p>
		{/snippet}
        {#snippet controller()}
            <Select.Root type="single" name="ioMode" bind:value={modeBandwidth}>
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
		{#snippet content()}
				<div class="h-[200px] w-full resize overflow-visible">
					<ScatterChart
						x="completedAt"
						y="bandwidthBytes"
                        series={Object.values(
                            modeBandwidth === "read" ? readTmp :
                            modeBandwidth === "write" ? writeTmp :
                            trimTmp
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
									value={`${Number(formatCapacity(data.bandwidthBytes).value).toFixed(0)} ${formatCapacity(data.bandwidthBytes).unit}`}
								/>
								<Tooltip.Item label="Date" value={dayjs(data.completedAt).format('YYYY/MM/DD HH:mm')} />
							</Tooltip.List>
							</Tooltip.Root>
						</svelte:fragment>
					</ScatterChart>
				</div>
		{/snippet}
	</Template.Area>

	<Template.Area title="IO">
		{#snippet hint()}
			<p>IO Per Second</p>
		{/snippet}
        {#snippet controller()}
            <Select.Root type="single" name="ioMode" bind:value={modeIO}>
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
		{#snippet content()}
				<div class="h-[200px] w-full resize overflow-visible">
					<ScatterChart
						x="completedAt"
						y="ioPerSecond"
                        series={Object.values(
                            modeIO === "read" ? readTmp :
                            modeIO === "write" ? writeTmp :
                            trimTmp
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
									label="IO"
									value={`${Number(formatCapacity(data.ioPerSecond).value).toFixed(0)} ${formatCapacity(data.ioPerSecond).unit}/s`}
								/>
								<Tooltip.Item label="Date" value={dayjs(data.completedAt).format('YYYY/MM/DD HH:mm')} />
							</Tooltip.List>
							</Tooltip.Root>
						</svelte:fragment>
					</ScatterChart>
				</div>
		{/snippet}
	</Template.Area>

	<Template.Area title="Latency">
		{#snippet hint()}
			<p>Latency</p>
		{/snippet}
        {#snippet controller()}
            <Select.Root type="single" name="ioMode" bind:value={modeLatency}>
            <Select.Trigger class="w-[180px]">
                {triggerLatencyContent}
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
		{#snippet content()}
				<div class="h-[200px] w-full resize overflow-visible">
					<ScatterChart
						x="completedAt"
						y="latency"
                        series={Object.values(
                            modeLatency === "read" ? readTmp :
                            modeLatency === "write" ? writeTmp :
                            trimTmp
                        )}
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
									label="IO"
									value={`${Number(formatCapacity(data.latency).value).toFixed(0)} ${formatCapacity(data.latency).unit}`}
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