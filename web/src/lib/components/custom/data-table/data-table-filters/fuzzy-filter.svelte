<script lang="ts" generics="TData">
	import * as Command from '$lib/components/ui/command/index.js';
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';
	import { cn } from '$lib/utils';
	import { capitalizeFirstLetter } from 'better-auth';

	let { table, columnId, alias }: { table: Table<TData>; columnId: string; alias?: string } =
		$props();

	const suggestions = $derived(
		([
			...new Set(table.getCoreRowModel().rows.map((row) => row.getValue(columnId)))
		].sort() as string[]) ?? ([] as string[])
	);
	let suggestionsOpen = $state(false);
</script>

<div class="relative max-w-xs">
	<Command.Root>
		<Command.Input
			class="h-10 placeholder:text-xs"
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
	</Command.Root>
</div>
