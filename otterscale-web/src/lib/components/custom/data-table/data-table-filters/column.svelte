<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';
</script>

<script lang="ts" generics="TData">
	let { table }: { table: Table<TData> } = $props();
</script>

<Popover.Root>
	<Popover.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="outline" size="sm" class="ml-auto text-xs">
				<Icon icon="ph:sliders-horizontal" />
				View
			</Button>
		{/snippet}
	</Popover.Trigger>
	<Popover.Content align="end" class="w-fit p-0">
		<Command.Root>
			<Command.Input placeholder="Search" />
			<Command.List>
				<Command.Group>
					{#each table
						.getAllColumns()
						.filter((column) => column.getCanHide()) as column (column.id)}
						<Command.Item
							onSelect={() => column.toggleVisibility(!column.getIsVisible())}
							class="text-xs capitalize"
						>
							<Icon
								icon="ph:check"
								class={cn('h-4 w-4', column.getIsVisible() ? 'visible' : 'invisible')}
							/>
							{column.id}
						</Command.Item>
					{/each}
				</Command.Group>
				<Command.Separator />
				<Command.Item
					onSelect={() => table.toggleAllColumnsVisible(!table.getIsAllColumnsVisible())}
					class="flex w-full items-center justify-center text-xs font-bold"
				>
					{table.getIsAllColumnsVisible() ? 'Clear' : 'All'}
				</Command.Item>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
