<script lang="ts" generics="TData">
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import { Badge } from '$lib/components/ui/badge';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let {
		table,
		columnId,
		messages,
		values,
		descriptor = (v: any) => v
	}: {
		table: Table<TData>;
		columnId: string;
		messages: Record<string, string>;
		values: any[];
		descriptor?: (v: any) => string;
	} = $props();

	const options = $derived(([...new Set(values)].sort() as string[]) ?? ([] as string[]));
	const extractions = $derived(
		(table.getColumn(columnId)?.getFilterValue() as string[]) ?? ([] as string[])
	);
</script>

<Popover.Root>
	<Popover.Trigger
		class={cn(buttonVariants({ size: 'sm', variant: 'outline' }), 'text-xs capitalize')}
	>
		<Icon icon="ph:funnel" />
		{messages[columnId]}

		{#if extractions.length > 0}
			<Separator orientation="vertical" />
		{/if}

		{#if extractions.length === 1}
			{@const [filteredValue] = extractions}
			<Badge variant="outline" class="-my-1 rounded-lg text-xs">
				{String(descriptor(filteredValue))}
			</Badge>
		{:else if extractions.length > 1}
			<HoverCard.Root>
				<HoverCard.Trigger>
					<span class="flex items-center gap-1">
						{extractions.length}
						<Icon icon="ph:checks" />
					</span>
				</HoverCard.Trigger>
				<HoverCard.Content class="flex w-fit flex-col gap-2 p-2">
					{#each extractions as filter}
						<Badge variant="outline" class="flex items-center gap-1 rounded-lg text-xs">
							<Icon icon="ph:funnel-simple" />
							{descriptor(filter)}
						</Badge>
					{/each}
				</HoverCard.Content>
			</HoverCard.Root>
		{/if}
	</Popover.Trigger>
	<Popover.Content class="w-fit p-0">
		<Command.Root>
			<Command.Input placeholder="Search" class="placeholder:text-xs" />
			<Command.List>
				<Command.Empty>Not found</Command.Empty>
				<Command.Group>
					{#each options as option}
						{@const newValue = extractions.includes(option)
							? extractions.filter((v) => v !== option)
							: [...extractions, option]}
						<Command.Item
							value={option}
							onSelect={() => {
								table.getColumn(columnId)?.setFilterValue(newValue.length ? newValue : undefined);
								table.firstPage();
							}}
						>
							<div class="flex w-full items-center gap-1 text-xs">
								<Icon
									icon={extractions.includes(option) ? 'ph:check' : 'ph:funnel-simple'}
									class={cn('h-4 w-4')}
								/>
								{descriptor(option)}
							</div>
						</Command.Item>
					{/each}
				</Command.Group>
				{#if extractions.length > 0}
					<Command.Separator />
					<Command.Item
						onSelect={() => table.getColumn(columnId)?.setFilterValue(undefined)}
						class="items-center justify-center text-xs font-bold hover:cursor-pointer"
					>
						{m.clear()}
					</Command.Item>
				{/if}
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
