<script lang="ts" module>
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';
	import { capitalizeFirstLetter } from 'better-auth';
</script>

<script lang="ts" generics="TData">
	let { table, columnId, alias }: { table: Table<TData>; columnId: string; alias?: string } =
		$props();

	let selectedValue: boolean | undefined = $state(undefined);

	const values = [false, true];
	const distinctValueCounts = $derived(
		values.reduce(
			(a, value) => {
				a[String(value)] = table
					.getCoreRowModel()
					.rows.filter((row) => row.getValue(columnId) === value).length;
				return a;
			},
			{} as Record<string, number>
		)
	);
</script>

<Popover.Root>
	<Popover.Trigger>
		<Button variant="outline" size="sm">
			<span class="flex items-center gap-1">
				<Icon icon="ph:funnel" />
				{alias ?? capitalizeFirstLetter(columnId)}
			</span>
			{#if selectedValue !== undefined}
				<Badge variant="outline">{selectedValue}</Badge>
			{/if}
		</Button>
	</Popover.Trigger>
	<Popover.Content class="w-fit p-0">
		<Command.Root>
			<Command.List>
				<Command.Group>
					{#each values as value}
						<Command.Item
							value={String(value)}
							onclick={() => {
								selectedValue = value;
								table.getColumn(columnId)?.setFilterValue(value);
							}}
							class="borderr flex w-full items-center justify-between text-xs hover:cursor-pointer"
						>
							<div class="flex w-full items-center gap-1 text-xs">
								<Icon
									icon={selectedValue === value ? 'ph:check' : 'ph:funnel-simple'}
									class={cn('h-4 w-4')}
								/>
								{capitalizeFirstLetter(String(value))}
								<p class="text-muted-foreground ml-auto font-mono">
									{distinctValueCounts[String(value)]}
								</p>
							</div>
						</Command.Item>
					{/each}
				</Command.Group>
				<Command.Separator />
				<Command.Item
					onclick={() => {
						table.getColumn(columnId)?.setFilterValue(undefined);
						selectedValue = undefined;
					}}
					class="w-full items-center justify-center text-xs font-bold hover:cursor-pointer"
					>Clear
				</Command.Item>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
