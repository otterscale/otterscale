<script lang="ts" module>
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';
</script>

<script lang="ts" generics="TData">
	let { table, messages }: { table: Table<TData>; messages: Record<string, string> } = $props();
</script>

<Popover.Root>
	<Popover.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="outline" size="sm" class="ml-auto text-xs">
				<Icon icon="ph:sliders-horizontal" />
				{m.datatable_filter_columns()}
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
							class="text-xs uppercase"
						>
							<Icon
								icon="ph:check"
								class={cn('h-4 w-4', column.getIsVisible() ? 'visible' : 'invisible')}
							/>
							{messages[column.id]}
						</Command.Item>
					{/each}
				</Command.Group>
				<Command.Separator />
				<Command.Item
					onSelect={() => table.toggleAllColumnsVisible(!table.getIsAllColumnsVisible())}
					class="flex w-full items-center justify-center text-xs font-bold"
				>
					{table.getIsAllColumnsVisible()
						? m.datatable_filter_action_clear()
						: m.datatable_filter_action_all()}
				</Command.Item>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
