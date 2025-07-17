<script lang="ts" module>
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';
	import { capitalizeFirstLetter } from 'better-auth';
</script>

<script lang="ts" generics="TData">
	let { table, columnId, alias }: { table: Table<TData>; columnId: string; alias?: string } =
		$props();

	const values = $derived(
		([
			...new Set(table.getCoreRowModel().rows.flatMap((row) => row.getValue(columnId)))
		].sort() as string[]) ?? ([] as string[])
	);
	const filteredValues = $derived(
		(table.getColumn(columnId)?.getFilterValue() as string[]) ?? ([] as string[])
	);
	const distinctValueCounts = $derived(
		values.reduce(
			(a, value) => {
				a[value] = table
					.getCoreRowModel()
					.rows.filter((row) => (row.getValue(columnId) as string[]).includes(value)).length;
				return a;
			},
			{} as Record<string, number>
		)
	);
</script>

<Popover.Root>
	<Popover.Trigger>
		{@render filterTrigger()}
	</Popover.Trigger>
	<Popover.Content class="w-fit p-0">
		<Command.Root>
			<Command.Input placeholder="Search" class="placeholder:text-xs" />
			<Command.List>
				<Command.Empty>Not found</Command.Empty>
				<Command.Group>
					{@render filterItems()}
				</Command.Group>
				{@render filterActions()}
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>

{#snippet filterTrigger()}
	<Button variant="outline" size="sm">
		<span class="flex items-center gap-1">
			<Icon icon="ph:funnel" />
			{alias ?? capitalizeFirstLetter(columnId)}
		</span>

		{#if filteredValues.length > 0}
			<Separator orientation="vertical" />
		{/if}

		{#if filteredValues.length === 1}
			{@const [filteredValue] = filteredValues}
			{@render showFilter(filteredValue)}
		{:else if filteredValues.length > 1}
			{@render listFilters(filteredValues)}
		{/if}
	</Button>
{/snippet}

{#snippet filterItems()}
	{#each values as value}
		{@const newValue = filteredValues.includes(value)
			? filteredValues.filter((v) => v !== value)
			: [...filteredValues, value]}
		<Command.Item
			{value}
			onSelect={() => {
				table.getColumn(columnId)?.setFilterValue(newValue.length ? newValue : undefined);
			}}
		>
			<div class="flex w-full items-center gap-1 text-xs">
				<Icon
					icon={filteredValues.includes(value) ? 'ph:check' : 'ph:funnel-simple'}
					class={cn('h-4 w-4')}
				/>
				{value}
				<p class="text-muted-foreground ml-auto font-mono">{distinctValueCounts[value]}</p>
			</div>
		</Command.Item>
	{/each}
{/snippet}

{#snippet filterActions()}
	{#if filteredValues.length > 0}
		<Command.Separator />
		<Command.Item
			onSelect={() => table.getColumn(columnId)?.setFilterValue(undefined)}
			class="items-center justify-center text-xs font-bold hover:cursor-pointer"
		>
			Clear
		</Command.Item>
	{/if}
{/snippet}

{#snippet showFilter(filter: any)}
	<Badge variant="secondary" class="-my-1 rounded-lg text-xs">
		{String(filter)}
	</Badge>
{/snippet}

{#snippet listFilters(filteredValues: string[])}
	<HoverCard.Root>
		<HoverCard.Trigger>
			<span class="flex items-center gap-1">
				{filteredValues.length}
				<Icon icon="ph:checks" />
			</span>
		</HoverCard.Trigger>
		<HoverCard.Content class="flex w-fit flex-col gap-2 p-2">
			{#each filteredValues as filter}
				<Badge variant="outline" class="flex items-center gap-1 rounded-lg text-xs">
					<Icon icon="ph:funnel-simple" />
					{filter}
				</Badge>
			{/each}
		</HoverCard.Content>
	</HoverCard.Root>
{/snippet}
