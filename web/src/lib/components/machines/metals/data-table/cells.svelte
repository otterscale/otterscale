<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
	import { AreaChart } from 'layerchart';
	import { SampleValue } from 'prometheus-query';

	import { resolve } from '$app/paths';
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import * as Chart from '$lib/components/ui/chart';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';
	import { cn } from '$lib/utils';

	import type { Metrics } from '../types';
	import Actions from './cell-actions.svelte';
	import GPUs from './cell-gpus.svelte';
	import Tags from './cell-tags.svelte';

	export const cells = {
		row_picker,
		fqdn_ip,
		power_state,
		status,
		cores_arch,
		ram,
		disk,
		storage,
		gpu,
		tags,
		scope,
		cpu_metric,
		memory_metric,
		storage_metric,
		actions
	};
</script>

{#snippet row_picker(row: Row<Machine>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet fqdn_ip(row: Row<Machine>)}
	<Table.Cell alignClass="items-start">
		<a
			class="m-0 p-0 underline hover:no-underline"
			href={resolve('/(auth)/machines/metal/[id]', {
				id: row.original.id
			})}
		>
			{row.original.fqdn}
		</a>
		{#if row.original.ipAddresses}
			<Table.SubCell>
				{#each row.original.ipAddresses as ipAddress}
					{ipAddress}
				{/each}
			</Table.SubCell>
		{/if}
	</Table.Cell>
{/snippet}

{#snippet power_state(row: Row<Machine>)}
	<div class="flex items-center gap-1">
		<Icon
			icon={row.original.powerState === 'on' ? 'ph:power' : 'ph:power'}
			class={cn(
				'size-4',
				row.original.powerState === 'on' ? 'text-accent-foreground' : 'text-destructive'
			)}
		/>
		<Table.Cell alignClass="items-start">
			{row.original.powerState}
			<Table.SubCell>
				{row.original.powerType}
			</Table.SubCell>
		</Table.Cell>
	</div>
{/snippet}

{#snippet status(row: Row<Machine>)}
	<Table.Cell alignClass="items-start">
		{@const processingStates = [
			'commissioning',
			'deploying',
			'disk_erasing',
			'entering_rescue_mode',
			'exiting_rescue_mode',
			'releasing',
			'testing'
		]}
		<Badge variant="outline">
			{row.original.status}
		</Badge>
		<Table.SubCell>
			{#if row.original.statusMessage != 'Deployed'}
				<span class="flex items-center gap-1">
					{#if processingStates.includes(row.original.status.toLowerCase())}
						<Icon icon="ph:spinner" class="animate-spin" />
					{/if}
					<p class="invisible max-w-[300px] truncate lg:visible">
						{row.original.statusMessage}
					</p>
				</span>
			{:else}
				<p class="invisible lg:visible">
					{`${row.original.osystem} ${row.original.hweKernel} ${row.original.distroSeries}`}
				</p>
			{/if}
		</Table.SubCell>
	</Table.Cell>
{/snippet}

{#snippet cores_arch(row: Row<Machine>)}
	<Table.Cell alignClass="items-right">
		{row.original.cpuCount}
		<Table.SubCell>
			{row.original.architecture}
		</Table.SubCell>
	</Table.Cell>
{/snippet}

{#snippet ram(row: Row<Machine>)}
	{@const { value, unit } = formatCapacity(Number(row.original.memoryMb) * 1000 * 1000)}
	<Table.Cell alignClass="items-end">
		{value}
		{unit}
	</Table.Cell>
{/snippet}

{#snippet disk(row: Row<Machine>)}
	<Table.Cell alignClass="items-end">
		{row.original.blockDevices.length}
	</Table.Cell>
{/snippet}

{#snippet storage(row: Row<Machine>)}
	{@const { value, unit } = formatCapacity(Number(row.original.storageMb) * 1000 * 1000)}
	<Table.Cell alignClass="items-end">
		{value}
		{unit}
	</Table.Cell>
{/snippet}

<!-- TODO: fix scope -->
{#snippet gpu(row: Row<Machine>)}
	<Table.Cell alignClass="items-end">
		<GPUs machine={row.original} />
	</Table.Cell>
{/snippet}

{#snippet tags(data: { row: Row<Machine>; reloadManager: ReloadManager })}
	<Table.Cell alignClass="items-start">
		<Tags machine={data.row.original} reloadManager={data.reloadManager} />
	</Table.Cell>
{/snippet}

{#snippet scope(row: Row<Machine>)}
	{@const identifier = row.original.workloadAnnotations['juju-machine-id']}
	<Table.Cell alignClass="items-start">
		{#if identifier}
			{@const scope = identifier.split('-machine-')[0]}
			{scope}
			{#if row.original.lastCommissioned}
				<Table.SubCell>
					{formatTimeAgo(timestampDate(row.original.lastCommissioned))}
				</Table.SubCell>
			{/if}
		{/if}
	</Table.Cell>
{/snippet}

{#snippet cpu_metric(data: { row: Row<Machine>; metrics: Metrics })}
	{@const usages: SampleValue[] = data.metrics?.cpu.get(data.row.original.fqdn) ?? []}
	{@const maximumValue = Math.max(...usages.map((usage) => Number(usage.value)))}
	{@const minimumValue = Math.min(...usages.map((usage) => Number(usage.value)))}
	{@const configuration = {
		value: { label: 'usage', color: maximumValue > 0.5 ? 'var(--warning)' : 'var(--healthy)' }
	} satisfies Chart.ChartConfig}
	{#if usages.length > 0}
		<Table.Cell alignClass="relative justify-center">
			<Chart.Container config={configuration} class="h-10 w-full">
				<AreaChart
					data={usages}
					x="time"
					series={[
						{
							key: 'value',
							label: configuration['value'].label,
							color: configuration['value'].color
						}
					]}
					props={{
						area: {
							curve: curveNatural,
							'fill-opacity': 0.1,
							line: { class: 'stroke-1' },
							motion: 'tween'
						}
					}}
					axis={false}
					xScale={scaleUtc()}
					yDomain={[minimumValue, maximumValue]}
					grid={false}
				>
					{#snippet tooltip()}
						<Chart.Tooltip
							indicator="dot"
							labelFormatter={(v: Date) => {
								return v.toLocaleDateString('en-US', {
									year: 'numeric',
									month: 'short',
									day: 'numeric',
									hour: 'numeric',
									minute: 'numeric'
								});
							}}
						>
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
				</AreaChart>
			</Chart.Container>
		</Table.Cell>
	{/if}
{/snippet}

{#snippet memory_metric(data: { row: Row<Machine>; metrics: Metrics })}
	{@const usages: SampleValue[] = data.metrics?.memory.get(data.row.original.fqdn) ?? []}
	{@const maximumValue = Math.max(...usages.map((usage) => Number(usage.value)))}
	{@const minimumValue = Math.min(...usages.map((usage) => Number(usage.value)))}
	{@const configuration = {
		value: { label: 'usage', color: 'var(--chart-3)' }
	} satisfies Chart.ChartConfig}
	{#if usages.length > 0}
		<Table.Cell alignClass="relative justify-center">
			<Chart.Container config={configuration} class="h-10 w-full">
				<AreaChart
					data={usages}
					x="time"
					series={[
						{
							key: 'value',
							label: configuration['value'].label,
							color: configuration['value'].color
						}
					]}
					props={{
						area: {
							curve: curveNatural,
							'fill-opacity': 0.1,
							line: { class: 'stroke-1' },
							motion: 'tween'
						}
					}}
					axis={false}
					yDomain={[minimumValue, maximumValue]}
					xScale={scaleUtc()}
					grid={false}
				>
					{#snippet tooltip()}
						<Chart.Tooltip
							indicator="dot"
							labelFormatter={(v: Date) => {
								return v.toLocaleDateString('en-US', {
									year: 'numeric',
									month: 'short',
									day: 'numeric',
									hour: 'numeric',
									minute: 'numeric'
								});
							}}
						>
							{#snippet formatter({ item, name, value })}
								{@const { value: usageValue, unit: usageUnit } = formatCapacity(Number(value))}
								<div
									class="flex flex-1 shrink-0 items-center justify-start gap-1 font-mono text-xs leading-none"
									style="--color-bg: {item.color}"
								>
									<Icon icon="ph:square-fill" class="text-(--color-bg)" />
									<h1 class="font-semibold text-muted-foreground">{name}</h1>
									<p class="ml-auto">{usageValue} {usageUnit}</p>
								</div>
							{/snippet}
						</Chart.Tooltip>
					{/snippet}
				</AreaChart>
			</Chart.Container>
		</Table.Cell>
	{/if}
{/snippet}

{#snippet storage_metric(data: { row: Row<Machine>; metrics: Metrics })}
	{@const usages: SampleValue[] = data.metrics?.storage.get(data.row.original.fqdn) ?? []}
	{@const maximumValue = Math.max(...usages.map((usage) => Number(usage.value)))}
	{@const minimumValue = Math.min(...usages.map((usage) => Number(usage.value)))}
	{@const configuration = {
		value: { label: 'usage', color: maximumValue > 0.5 ? 'var(--warning)' : 'var(--healthy)' }
	} satisfies Chart.ChartConfig}
	{#if usages.length > 0}
		<Table.Cell alignClass="relative justify-center">
			<Chart.Container config={configuration} class="h-10 w-full">
				<AreaChart
					data={usages}
					x="time"
					series={[
						{
							key: 'value',
							label: configuration['value'].label,
							color: configuration['value'].color
						}
					]}
					props={{
						area: {
							curve: curveNatural,
							'fill-opacity': 0.1,
							line: { class: 'stroke-1' },
							motion: 'tween'
						}
					}}
					axis={false}
					xScale={scaleUtc()}
					yDomain={[minimumValue, maximumValue]}
					grid={false}
				>
					{#snippet tooltip()}
						<Chart.Tooltip
							indicator="dot"
							labelFormatter={(v: Date) => {
								return v.toLocaleDateString('en-US', {
									year: 'numeric',
									month: 'short',
									day: 'numeric',
									hour: 'numeric',
									minute: 'numeric'
								});
							}}
						>
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
				</AreaChart>
			</Chart.Container>
		</Table.Cell>
	{/if}
{/snippet}

{#snippet actions(data: { row: Row<Machine>; reloadManager: ReloadManager })}
	<Table.Cell alignClass="items-start">
		<Actions machine={data.row.original} reloadManager={data.reloadManager} />
	</Table.Cell>
{/snippet}
