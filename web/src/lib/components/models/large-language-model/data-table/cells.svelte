<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';
	import { scaleUtc } from 'd3-scale';
	import { curveLinear } from 'd3-shape';
	import { LineChart } from 'layerchart';

	import { type LargeLanguageModel } from '../type';

	import Actions from './cell-actions.svelte';
	import Relation from './cell-relation.svelte';

	import { page } from '$app/state';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as Chart from '$lib/components/ui/chart';
	import { dynamicPaths } from '$lib/path';

	export const cells = {
		row_picker,
		model,
		name,
		replicas,
		healthies,
		gpu_cache,
		kv_cache,
		requests,
		time_to_first_token,
		relation,
		action,
	};
</script>

{#snippet row_picker(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet model(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-start">
		<a
			class="m-0 p-0 underline hover:no-underline"
			href={`${dynamicPaths.applicationsWorkloads(page.params.scope).url}/${row.original.application.namespace}/${row.original.application.name}`}
		>
			{row.original.application.name}
		</a>
		<Layout.SubCell>
			{row.original.application.namespace}
		</Layout.SubCell>
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet replicas(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		{row.original.application.replicas}
	</Layout.Cell>
{/snippet}

{#snippet healthies(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		{row.original.application.healthies}
	</Layout.Cell>
{/snippet}

{#snippet gpu_cache(row: Row<LargeLanguageModel>)}
	{@const configuration = {
		request: { label: 'request', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig}

	<Layout.Cell class="items-end">
		<Chart.Container config={configuration} class="h-fit w-20">
			<LineChart
				data={row.original.metrics.gpu_cache}
				x="time"
				xScale={scaleUtc()}
				axis={false}
				series={[
					{
						key: 'value',
						label: configuration.request.label,
						color: configuration.request.color,
					},
				]}
				props={{
					spline: { curve: curveLinear, motion: 'tween', strokeWidth: 2 },
					xAxis: {
						format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' }),
					},
					highlight: { points: { r: 4 } },
				}}
			>
				{#snippet tooltip()}
					<Chart.Tooltip hideLabel>
						{#snippet formatter({ item, name, value })}
							<div
								style="--color-bg: {item.color}"
								class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
							></div>
							<div class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none">
								<div class="grid gap-1.5">
									<span class="text-muted-foreground">{name}</span>
								</div>
								<p class="font-mono">{Number(value)}</p>
							</div>
						{/snippet}
					</Chart.Tooltip>
				{/snippet}
			</LineChart>
		</Chart.Container>
	</Layout.Cell>
{/snippet}

{#snippet kv_cache(row: Row<LargeLanguageModel>)}
	{@const configuration = {
		request: { label: 'request', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig}

	<Layout.Cell class="items-end">
		<Chart.Container config={configuration} class="h-fit w-20">
			<LineChart
				data={row.original.metrics.kv_cache}
				x="time"
				xScale={scaleUtc()}
				axis={false}
				series={[
					{
						key: 'value',
						label: configuration.request.label,
						color: configuration.request.color,
					},
				]}
				props={{
					spline: { curve: curveLinear, motion: 'tween', strokeWidth: 2 },
					xAxis: {
						format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' }),
					},
					highlight: { points: { r: 4 } },
				}}
			>
				{#snippet tooltip()}
					<Chart.Tooltip hideLabel>
						{#snippet formatter({ item, name, value })}
							<div
								style="--color-bg: {item.color}"
								class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
							></div>
							<div class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none">
								<div class="grid gap-1.5">
									<span class="text-muted-foreground">{name}</span>
								</div>
								<p class="font-mono">{Number(value)}</p>
							</div>
						{/snippet}
					</Chart.Tooltip>
				{/snippet}
			</LineChart>
		</Chart.Container>
	</Layout.Cell>
{/snippet}

{#snippet requests(row: Row<LargeLanguageModel>)}
	{@const configuration = {
		request: { label: 'request', color: 'var(--chart-1)' },
	} satisfies Chart.ChartConfig}

	<Layout.Cell class="items-end">
		<Chart.Container config={configuration} class="h-fit w-20">
			<LineChart
				data={row.original.metrics.requests}
				x="time"
				xScale={scaleUtc()}
				axis={false}
				series={[
					{
						key: 'value',
						label: configuration.request.label,
						color: configuration.request.color,
					},
				]}
				props={{
					spline: { curve: curveLinear, motion: 'tween', strokeWidth: 2 },
					xAxis: {
						format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' }),
					},
					highlight: { points: { r: 4 } },
				}}
			>
				{#snippet tooltip()}
					<Chart.Tooltip hideLabel>
						{#snippet formatter({ item, name, value })}
							<div
								style="--color-bg: {item.color}"
								class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
							></div>
							<div class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none">
								<div class="grid gap-1.5">
									<span class="text-muted-foreground">{name}</span>
								</div>
								<p class="font-mono">{Number(value)}</p>
							</div>
						{/snippet}
					</Chart.Tooltip>
				{/snippet}
			</LineChart>
		</Chart.Container>
	</Layout.Cell>
{/snippet}

{#snippet time_to_first_token(row: Row<LargeLanguageModel>)}
	{@const configuration = {
		request: { label: 'request', color: 'var(--chart-1)' },
	} satisfies Chart.ChartConfig}

	<Layout.Cell class="items-end">
		<Chart.Container config={configuration} class="h-fit w-20">
			<LineChart
				data={row.original.metrics.time_to_first_token}
				x="time"
				xScale={scaleUtc()}
				axis={false}
				series={[
					{
						key: 'value',
						label: configuration.request.label,
						color: configuration.request.color,
					},
				]}
				props={{
					spline: { curve: curveLinear, motion: 'tween', strokeWidth: 2 },
					xAxis: {
						format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' }),
					},
					highlight: { points: { r: 4 } },
				}}
			>
				{#snippet tooltip()}
					<Chart.Tooltip hideLabel>
						{#snippet formatter({ item, name, value })}
							<div
								style="--color-bg: {item.color}"
								class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
							></div>
							<div class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none">
								<div class="grid gap-1.5">
									<span class="text-muted-foreground">{name}</span>
								</div>
								<p class="font-mono">{Number(value)}</p>
							</div>
						{/snippet}
					</Chart.Tooltip>
				{/snippet}
			</LineChart>
		</Chart.Container>
	</Layout.Cell>
{/snippet}

{#snippet relation(row: Row<LargeLanguageModel>)}
	{#if row.original.application.healthies > 0}
		<Layout.Cell class="items-end">
			<Relation model={row.original} />
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet action(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		<Actions model={row.original} />
	</Layout.Cell>
{/snippet}
