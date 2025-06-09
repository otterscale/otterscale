<script lang="ts" generics="TData">
	import * as Command from '$lib/components/ui/command/index.js';
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button';

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

<div class="max-w-xs">
	<Command.Root>
		<div class="relative">
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
			{#if table.getColumn(columnId)?.getFilterValue()}
				<Button
					class="absolute right-2 top-1/2 -translate-y-1/2 p-1 text-gray-500 hover:text-gray-700"
					variant="ghost"
					size="icon"
					onclick={() => {
						table.getColumn(columnId)?.setFilterValue(undefined);
						suggestionsOpen = false;
					}}
				>
					<Icon icon="ph:x-circle" />
				</Button>
			{/if}
		</div>
		<Command.List
			class={cn(
				'absolute z-50 mt-10 w-fit rounded-md border bg-white shadow',
				suggestionsOpen ? 'visible' : 'hidden'
			)}
		>
			{#each suggestions as suggestion}
				<Command.Item
					disabled={!suggestionsOpen}
					value={suggestion}
					class="p-2 text-xs hover:cursor-pointer"
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
