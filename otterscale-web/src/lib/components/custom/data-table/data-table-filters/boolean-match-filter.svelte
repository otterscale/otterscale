<script lang="ts" module>
	import { Badge } from '$lib/components/ui/badge';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
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
		values: boolean[];
		descriptor?: (v: any) => string;
	} = $props();

	let selectedValue: boolean | undefined = $state(undefined);

	const options = [false, true];
	const counts = $derived(
		options.reduce(
			(a, option) => {
				a[String(option)] = values.filter((value) => value === option).length;
				return a;
			},
			{} as Record<string, number>
		)
	);
</script>

<Popover.Root>
	<Popover.Trigger class={cn(buttonVariants({ size: 'sm', variant: 'outline' }))}>
		<Icon icon="ph:funnel" />
		{alias ?? capitalizeFirstLetter(columnId)}
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
							}}
							class="flex w-full items-center justify-between gap-4 text-xs"
						>
							<span class="flex w-full items-center gap-1 text-xs">
								<Icon
									icon={selectedValue === option ? 'ph:check' : 'ph:funnel-simple'}
									class={cn('h-4 w-4')}
								/>
								{capitalizeFirstLetter(String(descriptor(option)))}
							</span>
							<p class="text-muted-foreground ml-auto font-mono">
								{counts[String(option)]}
							</p>
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
					Clear
				</Command.Item>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
