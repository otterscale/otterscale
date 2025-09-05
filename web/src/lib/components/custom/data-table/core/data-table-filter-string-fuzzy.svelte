<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import { buttonVariants } from '$lib/components/ui/button';
	import * as Command from '$lib/components/ui/command/index.js';
	import { cn } from '$lib/utils';
</script>

<script lang="ts" generics="TData">
	let {
		table,
		columnId,
		messages,
		values,
	}: { table: Table<TData>; columnId: string; messages: Record<string, string>; values: any[] } = $props();

	const suggestions = $derived(([...new Set(values)].sort() as string[]) ?? ([] as string[]));
	let open = $state(false);
</script>

<div class="relative w-40">
	<Command.Root class={cn(buttonVariants({ size: 'sm', variant: 'outline' }), 'text-xs')}>
		<div class="relative">
			<Command.Input
				class="pr-3 placeholder:text-xs placeholder:uppercase"
				placeholder={messages[columnId]}
				value={(table.getColumn(columnId)?.getFilterValue() as string) ?? ''}
				oninput={(e) => {
					table.getColumn(columnId)?.setFilterValue(e.currentTarget.value);
					table.firstPage();
					open = e.currentTarget.value ? true : false;
				}}
				onmousedowncapture={() => {
					open = true;
				}}
				onblur={() => {
					open = false;
				}}
			/>
			<button
				class={cn(
					'absolute top-1/2 right-0 -translate-y-1/2',
					table.getColumn(columnId)?.getFilterValue() ? 'visible' : 'hidden',
				)}
				onclick={() => {
					table.getColumn(columnId)?.setFilterValue(undefined);
				}}
			>
				<Icon icon="ph:x" />
			</button>
		</div>
		<Command.List
			class={cn(
				'bg-card absolute top-10 left-0 z-50 w-fit min-w-40 rounded-md border shadow',
				open ? 'visible' : 'hidden',
			)}
		>
			{#each suggestions as suggestion}
				<Command.Item
					disabled={!open}
					value={suggestion}
					class="text-xs hover:cursor-pointer"
					onmousedown={(e) => {
						e.preventDefault();

						table.getColumn(columnId)?.setFilterValue(suggestion);
						table.firstPage();
						open = false;
					}}
				>
					<Icon icon="ph:list-magnifying-glass" />
					{suggestion}
				</Command.Item>
			{/each}
		</Command.List>
	</Command.Root>
</div>
