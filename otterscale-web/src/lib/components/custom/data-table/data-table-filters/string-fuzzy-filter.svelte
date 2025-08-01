<script lang="ts" module>
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Command from '$lib/components/ui/command/index.js';
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
		values
	}: { table: Table<TData>; columnId: string; alias?: string; values: any[] } = $props();

	const suggestions = $derived(([...new Set(values)].sort() as string[]) ?? ([] as string[]));
	let suggestionsOpen = $state(false);
</script>

<div class="relative max-w-xs">
	<Command.Root class={cn(buttonVariants({ size: 'sm', variant: 'outline' }), 'text-xs')}>
		{@render filterInput()}
		{@render filterSuggestions()}
	</Command.Root>
</div>

{#snippet filterInput()}
	<Command.Input
		class="placeholder:text-xs"
		placeholder={alias ?? capitalizeFirstLetter(columnId)}
		value={(table.getColumn(columnId)?.getFilterValue() as string) ?? ''}
		oninput={(e) => {
			table.getColumn(columnId)?.setFilterValue(e.currentTarget.value);
			suggestionsOpen = true;
			if (!e.currentTarget.value) {
				suggestionsOpen = false;
			}
		}}
		onmousedowncapture={() => {
			suggestionsOpen = true;
		}}
		onblur={(e) => {
			suggestionsOpen = false;
		}}
	/>
{/snippet}

{#snippet filterSuggestions()}
	<Command.List
		class={cn(
			'bg-card absolute top-10 z-50 w-full rounded-md border shadow',
			suggestionsOpen ? 'visible' : 'hidden'
		)}
	>
		{#each suggestions as suggestion}
			<Command.Item
				disabled={!suggestionsOpen}
				value={suggestion}
				class="text-xs hover:cursor-pointer"
				onmousedown={(e) => {
					e.preventDefault();
					table.getColumn(columnId)?.setFilterValue(suggestion);
					suggestionsOpen = false;
				}}
			>
				<Icon icon="ph:list-magnifying-glass" />
				{suggestion}
			</Command.Item>
		{/each}
	</Command.List>
{/snippet}
