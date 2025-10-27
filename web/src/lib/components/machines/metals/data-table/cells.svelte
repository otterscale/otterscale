<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';
	import GPUs from './cell-gpus.svelte';
	import Tags from './cell-tags.svelte';

	import { page } from '$app/state';
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';
	import { dynamicPaths } from '$lib/path';
	import { cn } from '$lib/utils';

	export const cells = {
		row_picker,
		fqdn_ip,
		powerState,
		status,
		cores_arch,
		ram,
		disk,
		storage,
		gpu,
		tags,
		scope,
		actions,
	};
</script>

{#snippet row_picker(row: Row<Machine>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet fqdn_ip(row: Row<Machine>)}
	<Layout.Cell class="items-start">
		<a
			class="m-0 p-0 underline hover:no-underline"
			href={`${dynamicPaths.machinesMetal(page.params.scope).url}/${row.original.id}`}
		>
			{row.original.fqdn}
		</a>
		{#if row.original.ipAddresses}
			<Layout.SubCell>
				{#each row.original.ipAddresses as ipAddress}
					{ipAddress}
				{/each}
			</Layout.SubCell>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet powerState(row: Row<Machine>)}
	<Layout.Cell class="flex-row items-center">
		<Icon
			icon={row.original.powerState === 'on' ? 'ph:power' : 'ph:power'}
			class={cn('size-4', row.original.powerState === 'on' ? 'text-accent-foreground' : 'text-destructive')}
		/>
		<Layout.Cell>
			{row.original.powerState}
			<Layout.SubCell>
				{row.original.powerType}
			</Layout.SubCell>
		</Layout.Cell>
	</Layout.Cell>
{/snippet}

{#snippet status(row: Row<Machine>)}
	<Layout.Cell class="items-start">
		{@const processingStates = [
			'commissioning',
			'deploying',
			'disk_erasing',
			'entering_rescue_mode',
			'exiting_rescue_mode',
			'releasing',
			'testing',
		]}
		<Badge variant="outline">
			{row.original.status}
		</Badge>
		<Layout.SubCell>
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
		</Layout.SubCell>
	</Layout.Cell>
{/snippet}

{#snippet cores_arch(row: Row<Machine>)}
	<Layout.Cell class="items-right">
		{row.original.cpuCount}
		<Layout.SubCell>
			{row.original.architecture}
		</Layout.SubCell>
	</Layout.Cell>
{/snippet}

{#snippet ram(row: Row<Machine>)}
	{@const { value, unit } = formatCapacity(Number(row.original.memoryMb) * 1000 * 1000)}
	<Layout.Cell class="items-end">
		{value}
		{unit}
	</Layout.Cell>
{/snippet}

{#snippet disk(row: Row<Machine>)}
	<Layout.Cell class="items-end">
		{row.original.blockDevices.length}
	</Layout.Cell>
{/snippet}

{#snippet storage(row: Row<Machine>)}
	{@const { value, unit } = formatCapacity(Number(row.original.storageMb) * 1000 * 1000)}
	<Layout.Cell class="items-end">
		{value}
		{unit}
	</Layout.Cell>
{/snippet}

{#snippet gpu(row: Row<Machine>)}
	<Layout.Cell class="items-end">
		<GPUs machine={row.original} />
	</Layout.Cell>
{/snippet}

{#snippet tags(row: Row<Machine>)}
	<Layout.Cell class="items-start">
		<Tags machine={row.original} />
	</Layout.Cell>
{/snippet}

{#snippet scope(row: Row<Machine>)}
	{@const identifier = row.original.workloadAnnotations['juju-machine-id']}
	<Layout.Cell class="items-start">
		{#if identifier}
			{@const scope = identifier.split('-machine-')[0]}
			{scope}
			{#if row.original.lastCommissioned}
				<Layout.SubCell>
					{formatTimeAgo(timestampDate(row.original.lastCommissioned))}
				</Layout.SubCell>
			{/if}
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<Machine>)}
	<Layout.Cell class="items-start">
		<Actions machine={row.original} />
	</Layout.Cell>
{/snippet}
