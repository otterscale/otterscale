<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import { curveMonotoneX } from 'd3-shape';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { SampleValue } from 'prometheus-query';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import type { ObjectStorageDaemon } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as Progress from '$lib/components/custom/progress';
	import { Badge } from '$lib/components/ui/badge';
	import * as Chart from '$lib/components/ui/chart';
	import { formatCapacity, formatIO } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	import type { Metrics } from '../types';
	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		name,
		state,
		osdUp,
		osdIn,
		exists,
		deviceClass,
		machine,
		placementGroupCount,
		usage,
		iops,
		throughput,
		actions
	};
</script>

{#snippet row_picker(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet state(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="flex-row items-center">
		{#if row.original.in}
			<Badge variant="outline">{m.osd_in()}</Badge>
		{/if}
		{#if row.original.up}
			<Badge variant="outline">{m.osd_up()}</Badge>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet osdUp()}{/snippet}

{#snippet osdIn()}{/snippet}

{#snippet exists(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-start">
		{#if !row.original.exists}
			<Icon icon="ph:x" class="text-destructive" />
		{:else}
			<Icon icon="ph:circle" class="text-primary" />
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet machine(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-start">
		<div class="flex items-center gap-1">
			<Badge variant="outline">
				{row.original.machine?.hostname}
			</Badge>
			<Icon
				icon="ph:arrow-square-out"
				class="hover:cursor-pointer"
				onclick={() => {
					goto(
						resolve('/(auth)/machines/metal/[id]', {
							id: row.original.machine?.id ?? ''
						})
					);
				}}
			/>
		</div>
	</Layout.Cell>
{/snippet}

{#snippet deviceClass(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.deviceClass}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet placementGroupCount(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-end">
		{row.original.placementGroupCount}
	</Layout.Cell>
{/snippet}

{#snippet usage(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-end">
		<Progress.Root
			numerator={Number(row.original.usedBytes)}
			denominator={Number(row.original.sizeBytes)}
			target="STB"
		>
			{#snippet ratio({ numerator, denominator })}
				{Progress.formatRatio(numerator, denominator)}
			{/snippet}
			{#snippet detail({ numerator, denominator })}
				{@const { value: numeratorValue, unit: numeratorUnit } = formatCapacity(numerator)}
				{@const { value: denominatorValue, unit: denominatorUnit } = formatCapacity(denominator)}
				{numeratorValue}
				{numeratorUnit}/{denominatorValue}
				{denominatorUnit}
			{/snippet}
		</Progress.Root>
	</Layout.Cell>
{/snippet}

{#snippet iops(data: { row: Row<ObjectStorageDaemon>; metrics: Metrics })}
	{@const configuration = {
		input: { label: 'input', color: 'var(--chart-1)' },
		output: { label: 'output', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig}
	{@const inputs = data.metrics.input?.get(data.row.original.name) as SampleValue[]}
	{@const outputs = data.metrics.output?.get(data.row.original.name) as SampleValue[]}
	{@const ios = inputs.map((input, index) => ({
		time: input.time,
		input: input.value,
		output: outputs[index]?.value ?? 0
	}))}
	{@const maximumValue = Math.max(
		...inputs.map((input) => Number(input.value)),
		...outputs.map((output) => Number(output.value))
	)}
	{@const minimumValue = Math.min(
		...inputs.map((input) => Number(input.value)),
		...outputs.map((output) => Number(output.value))
	)}
	<Layout.Cell class="items-center">
		<Chart.Container config={configuration} class="relative h-20 w-full">
			<AreaChart
				data={ios}
				x="time"
				yDomain={[minimumValue, maximumValue]}
				series={[
					{
						key: 'input',
						label: configuration.input.label,
						color: configuration.input.color
					},
					{
						key: 'output',
						label: configuration.output.label,
						color: configuration.output.color
					}
				]}
				props={{
					area: {
						curve: curveMonotoneX,
						'fill-opacity': 0.4,
						line: { class: 'stroke-1' },
						motion: 'tween'
					},
					xAxis: { format: () => '' },
					yAxis: { format: () => '' }
				}}
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
							{@const { value: ioValue, unit: ioUnit } = formatIO(Number(value))}
							<div
								class="flex flex-1 shrink-0 items-center justify-start gap-1 font-mono text-xs leading-none"
								style="--color-bg: {item.color}"
							>
								<Icon icon="ph:square-fill" class="text-(--color-bg)" />
								<h1 class="font-semibold text-muted-foreground">{name}</h1>
								<p class="ml-auto">{ioValue} {ioUnit}</p>
							</div>
						{/snippet}
					</Chart.Tooltip>
				{/snippet}
				{#snippet marks({ series, getAreaProps })}
					{#each series as s, i (s.key)}
						<LinearGradient
							stops={[s.color ?? '', 'color-mix(in lch, ' + s.color + ' 10%, transparent)']}
							vertical
						>
							{#snippet children({ gradient })}
								<Area {...getAreaProps(s, i)} fill={gradient} />
							{/snippet}
						</LinearGradient>
					{/each}
				{/snippet}
			</AreaChart>
		</Chart.Container>
	</Layout.Cell>
{/snippet}

{#snippet throughput(data: { row: Row<ObjectStorageDaemon>; metrics: Metrics })}
	{@const configuration = {
		read: { label: 'read', color: 'var(--chart-1)' },
		write: { label: 'write', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig}
	{@const reads = data.metrics.read?.get(data.row.original.name) as SampleValue[]}
	{@const writes = data.metrics.write?.get(data.row.original.name) as SampleValue[]}
	{@const throughputs = reads.map((read, index) => ({
		time: read.time,
		read: read.value,
		write: writes[index]?.value ?? 0
	}))}
	{@const maximumValue = Math.max(
		...reads.map((read) => Number(read.value)),
		...writes.map((write) => Number(write.value))
	)}
	{@const minimumValue = Math.min(
		...reads.map((read) => Number(read.value)),
		...writes.map((write) => Number(write.value))
	)}
	<Layout.Cell class="items-center">
		<Chart.Container config={configuration} class="relative h-20 w-full">
			<AreaChart
				data={throughputs}
				x="time"
				yDomain={[minimumValue, maximumValue]}
				series={[
					{
						key: 'read',
						label: configuration.read.label,
						color: configuration.read.color
					},
					{
						key: 'write',
						label: configuration.write.label,
						color: configuration.write.color
					}
				]}
				props={{
					area: {
						curve: curveMonotoneX,
						'fill-opacity': 0.4,
						line: { class: 'stroke-1' },
						motion: 'tween'
					},
					xAxis: { format: () => '' },
					yAxis: { format: () => '' }
				}}
			>
				{#snippet tooltip()}
					<Chart.Tooltip
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
							{@const { value: ioValue, unit: ioUnit } = formatIO(Number(value))}
							<div
								class="flex flex-1 shrink-0 items-center justify-start gap-1 font-mono text-xs leading-none"
								style="--color-bg: {item.color}"
							>
								<Icon icon="ph:square-fill" class="text-(--color-bg)" />
								<h1 class="font-semibold text-muted-foreground">{name}</h1>
								<p class="ml-auto">{ioValue} {ioUnit}</p>
							</div>
						{/snippet}
					</Chart.Tooltip>
				{/snippet}
				{#snippet marks({ series, getAreaProps })}
					{#each series as s, i (s.key)}
						<LinearGradient
							stops={[s.color ?? '', 'color-mix(in lch, ' + s.color + ' 10%, transparent)']}
							vertical
						>
							{#snippet children({ gradient })}
								<Area {...getAreaProps(s, i)} fill={gradient} />
							{/snippet}
						</LinearGradient>
					{/each}
				{/snippet}
			</AreaChart>
		</Chart.Container>
	</Layout.Cell>
{/snippet}

{#snippet actions(data: { row: Row<ObjectStorageDaemon>; scope: string })}
	<Layout.Cell class="items-start">
		<Actions osd={data.row.original} scope={data.scope} />
	</Layout.Cell>
{/snippet}
