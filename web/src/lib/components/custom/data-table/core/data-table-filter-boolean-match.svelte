<script lang="ts" generics="TData">
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import { Badge } from '$lib/components/ui/badge';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let {
		table,
		columnId,
		messages,
		descriptor = (v: any) => v
	}: {
		table: Table<TData>;
		columnId: string;
		messages: Record<string, string>;
		descriptor?: (v: any) => string;
	} = $props();

	let selectedValue: boolean | undefined = $state(undefined);

	const options = [false, true];
</script>

<Popover.Root>
	<Popover.Trigger
		class={cn(buttonVariants({ size: 'sm', variant: 'outline' }), 'text-xs capitalize')}
	>
		<Icon icon="ph:funnel" />
		{messages[columnId]}
		{#if selectedValue !== undefined}
			<Separator orientation="vertical" />
			<Badge variant="outline">{descriptor(selectedValue)}</Badge>
		{/if}
	</Popover.Trigger>
	<Popover.Content class="w-[300px] p-0">
		<Command.Root>
			<Command.List>
				<Command.Group>
					{#each options as option}
						<Command.Item
							value={String(option)}
							onclick={() => {
								if (selectedValue === option) {
									selectedValue = undefined;
									table.getColumn(columnId)?.setFilterValue(undefined);
								} else {
									selectedValue = option;
									table.getColumn(columnId)?.setFilterValue(option);
								}
								table.firstPage();
							}}
							class="flex w-full items-center justify-between gap-4 text-xs"
						>
							<span class="flex w-full items-center gap-1 text-xs capitalize">
								<Icon
									icon={selectedValue === option ? 'ph:check' : 'ph:funnel-simple'}
									class={cn('h-4 w-4')}
								/>
								{String(descriptor(option))}
							</span>
						</Command.Item>
					{/each}
				</Command.Group>
				<Command.Separator />
				<Command.Item
					onclick={() => {
						table.getColumn(columnId)?.setFilterValue(undefined);
						selectedValue = undefined;
					}}
					class="w-full items-center justify-center text-xs font-bold"
				>
					{m.clear()}
				</Command.Item>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
