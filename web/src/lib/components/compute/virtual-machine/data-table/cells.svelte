<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import { scaleUtc } from 'd3-scale';
	import { LineChart } from 'layerchart';
	import { SampleValue } from 'prometheus-query';

	import { resolve } from '$app/paths';
	import type { VirtualMachine } from '$lib/api/instance/v1/instance_pb';
	import { Disk } from '$lib/components/compute/virtual-machine/disk';
	import { Port } from '$lib/components/compute/virtual-machine/port';
	import { getStatusInfo } from '$lib/components/compute/virtual-machine/units/type';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { Badge } from '$lib/components/ui/badge';
	import * as Chart from '$lib/components/ui/chart';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatTimeAgo } from '$lib/formatter';

	import type { Metrics } from '../types';
	import Actions from './cell-actions.svelte';
	import VNC from './cell-vnc.svelte';

	export const cells = {
		row_picker,
		name,
		status,
		namespace,
		machineId,
		instanceType,
		disk,
		port,
		createTime,
		vnc,
		cpu_metric,
		memory_metric,
		storage_metric,
		actions
	};
</script>

{#snippet row_picker(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		<div class="flex items-center gap-1">
			{row.original.name}
		</div>
	</Layout.Cell>
{/snippet}

{#snippet status(row: Row<VirtualMachine>)}
	{@const statusInfo = getStatusInfo(row.original.status)}
	<Layout.Cell class="items-start">
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger>
					<Icon icon={statusInfo.icon} class={`${statusInfo.color} h-5 w-5`} />
				</Tooltip.Trigger>
				<Tooltip.Content>
					{statusInfo.text}
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	</Layout.Cell>
{/snippet}

{#snippet namespace(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.namespace}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet machineId(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		{#if row.original.machineId}
			<a
				class="m-0 p-0 underline hover:no-underline"
				href={resolve('/(auth)/machines/metal/[id]', {
					id: row.original.machineId
				})}
			>
				{row.original.hostname}
			</a>
			<Layout.SubCell>
				{row.original.ipAddresses}
			</Layout.SubCell>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet instanceType(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		{#if row.original.instanceType}
			<div class="flex items-center gap-1">
				<Badge variant="outline">
					{row.original.instanceType.name}
				</Badge>
				<HoverCard.Root>
					<HoverCard.Trigger>
						<Icon icon="ph:info" />
					</HoverCard.Trigger>
					<HoverCard.Content class="min-w-[300px]">
						<Table.Root>
							<Table.Body class="text-xs">
								{#if row.original.instanceType.name}
									<Table.Row>
										<Table.Head class="text-left">Name</Table.Head>
										<Table.Cell>
											<Badge variant="outline">{row.original.instanceType.name}</Badge>
										</Table.Cell>
									</Table.Row>
								{/if}
								{#if row.original.instanceType.namespace}
									<Table.Row>
										<Table.Head class="text-left">Namespace</Table.Head>
										<Table.Cell>{row.original.instanceType.namespace}</Table.Cell>
									</Table.Row>
								{/if}
								{#if row.original.instanceType.cpuCores}
									<Table.Row>
										<Table.Head class="text-left">CPU Cores</Table.Head>
										<Table.Cell>{row.original.instanceType.cpuCores}</Table.Cell>
									</Table.Row>
								{/if}
								{#if row.original.instanceType.memoryBytes}
									<Table.Row>
										<Table.Head class="text-left">Memory</Table.Head>
										<Table.Cell
											>{Number(row.original.instanceType.memoryBytes) / 1024 ** 3} GB</Table.Cell
										>
									</Table.Row>
								{/if}
								{#if row.original.instanceType.clusterWide !== undefined}
									<Table.Row>
										<Table.Head class="text-left">Cluster Wide</Table.Head>
										<Table.Cell>
											<Badge variant="outline">{row.original.instanceType.clusterWide}</Badge>
										</Table.Cell>
									</Table.Row>
								{/if}
							</Table.Body>
						</Table.Root>
					</HoverCard.Content>
				</HoverCard.Root>
			</div>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet disk(data: { row: Row<VirtualMachine>; scope: string })}
	<Layout.Cell class="items-end">
		<Disk virtualMachine={data.row.original} scope={data.scope} />
	</Layout.Cell>
{/snippet}

{#snippet port(data: { row: Row<VirtualMachine>; scope: string; reloadManager: ReloadManager })}
	<Layout.Cell class="items-end">
		<Port
			virtualMachine={data.row.original}
			scope={data.scope}
			reloadManager={data.reloadManager}
		/>
	</Layout.Cell>
{/snippet}

{#snippet createTime(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		{#if row.original.createdAt}
			<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{formatTimeAgo(timestampDate(row.original.createdAt))}
					</Tooltip.Trigger>
					<Tooltip.Content>
						{timestampDate(row.original.createdAt)}
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet vnc(data: { row: Row<VirtualMachine>; scope: string })}
	<Layout.Cell class="items-end">
		{#if data.row.original.status === 'Running'}
			<VNC virtualMachine={data.row.original} scope={data.scope} />
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet cpu_metric(data: { row: Row<VirtualMachine>; metrics: Metrics })}
	{@const configuation = {
		value: { label: 'usage', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig}
	{@const usage: SampleValue[] = data.metrics.cpu.get(data.row.original.name) ?? []}
	{#if usage.length > 0}
		<Layout.Cell class="justify-center">
			<Chart.Container config={configuation} class="h-10 w-full">
				<LineChart
					data={usage}
					x="time"
					series={[
						{
							key: 'value',
							label: configuation['value'].label,
							color: configuation['value'].color
						}
					]}
					axis={false}
					xScale={scaleUtc()}
					yDomain={[0, 1]}
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
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet memory_metric(data: { row: Row<VirtualMachine>; metrics: Metrics })}
	{@const configuation = {
		value: { label: 'usage', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig}
	{@const usage: SampleValue[] = data.metrics.memory.get(data.row.original.name) ?? []}
	{#if usage.length > 0}
		<Layout.Cell class="justify-center">
			<Chart.Container config={configuation} class="h-10 w-full">
				<LineChart
					data={usage}
					x="time"
					series={[
						{
							key: 'value',
							label: configuation['value'].label,
							color: configuation['value'].color
						}
					]}
					axis={false}
					xScale={scaleUtc()}
					yDomain={[0, 1]}
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
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet storage_metric(data: { row: Row<VirtualMachine>; metrics: Metrics })}
	{@const configuation = {
		value: { label: 'usage', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig}
	{@const usage: SampleValue[] = data.metrics.storage.get(data.row.original.name) ?? []}
	{#if usage.length > 0}
		<Layout.Cell class="justify-center">
			<Chart.Container config={configuation} class="h-10 w-full">
				<LineChart
					data={usage}
					x="time"
					series={[
						{
							key: 'value',
							label: configuation['value'].label,
							color: configuation['value'].color
						}
					]}
					axis={false}
					xScale={scaleUtc()}
					yDomain={[0, 1]}
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
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet actions(data: { row: Row<VirtualMachine>; scope: string; reloadManager: ReloadManager })}
	<Layout.Cell class="items-start">
		<Actions
			virtualMachine={data.row.original}
			scope={data.scope}
			reloadManager={data.reloadManager}
		/>
	</Layout.Cell>
{/snippet}
