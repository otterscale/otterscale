<script lang="ts" module>
	import { Badge } from '$lib/components/ui/badge';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
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
	let {
		table,
		columnId,
		alias,
		values,
		descriptor = (v: any) => v
	}: {
		table: Table<TData>;
		columnId: string;
		alias?: string;
		values: any[];
		descriptor?: (v: any) => string;
	} = $props();

	const options = $derived(([...new Set(values)].sort() as string[]) ?? ([] as string[]));
	const extractions = $derived(
		(table.getColumn(columnId)?.getFilterValue() as string[]) ?? ([] as string[])
	);
	const counts = $derived(
		options.reduce(
			(a, option) => {
				a[option] = table
					.getCoreRowModel()
					.rows.filter((row) => row.getValue(columnId) === option).length;
				return a;
			},
			{} as Record<string, number>
		)
	);
</script>

<Popover.Root>
	<Popover.Trigger class={buttonVariants({ size: 'sm', variant: 'outline' })}>
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
	<Icon icon="ph:funnel" />
	{alias ?? capitalizeFirstLetter(columnId)}

	{#if extractions.length > 0}
		<Separator orientation="vertical" />
	{/if}

	{#if extractions.length === 1}
		{@const [filteredValue] = extractions}
		{@render showFilter(filteredValue)}
	{:else if extractions.length > 1}
		{@render listFilters(extractions)}
	{/if}
{/snippet}

{#snippet filterItems()}
	{#each options as option}
		{@const newValue = extractions.includes(option)
			? extractions.filter((v) => v !== option)
			: [...extractions, option]}
		<Command.Item
			value={option}
			onSelect={() => {
				table.getColumn(columnId)?.setFilterValue(newValue.length ? newValue : undefined);
			}}
		>
			<div class="flex w-full items-center gap-1 text-xs">
				<Icon
					icon={extractions.includes(option) ? 'ph:check' : 'ph:funnel-simple'}
					class={cn('h-4 w-4')}
				/>
				{descriptor(option)}
				<p class="text-muted-foreground ml-auto font-mono">{counts[option]}</p>
			</div>
		</Command.Item>
	{/each}
{/snippet}

{#snippet filterActions()}
	{#if extractions.length > 0}
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
	<Badge variant="outline" class="-my-1 rounded-lg text-xs">
		{String(descriptor(filter))}
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
					{descriptor(filter)}
				</Badge>
			{/each}
		</HoverCard.Content>
	</HoverCard.Root>
{/snippet}
