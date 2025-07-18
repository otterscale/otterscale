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
	const distinctValueCounts = $derived(
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
	<Popover.Trigger>
		<Button variant="outline" size="sm">
			<span class="flex items-center gap-1">
				<Icon icon="ph:funnel" />
				{alias ?? capitalizeFirstLetter(columnId)}
			</span>
			{#if selectedValue !== undefined}
				<Badge variant="outline">{descriptor(selectedValue)}</Badge>
			{/if}
		</Button>
	</Popover.Trigger>
	<Popover.Content class="w-[300px] p-0">
		<Command.Root>
			<Command.List>
				<Command.Group>
					{#each options as option}
						<Command.Item
							value={String(option)}
							onclick={() => {
								selectedValue = option;
								table.getColumn(columnId)?.setFilterValue(option);
							}}
							class="borderr flex w-full items-center justify-between text-xs hover:cursor-pointer"
						>
							<div class="flex w-full items-center gap-4 text-xs">
								<div class="flex w-full items-center gap-1 text-xs">
									<Icon
										icon={selectedValue === option ? 'ph:check' : 'ph:funnel-simple'}
										class={cn('h-4 w-4')}
									/>
									{capitalizeFirstLetter(String(descriptor(option)))}
								</div>
								<p class="text-muted-foreground ml-auto font-mono">
									{distinctValueCounts[String(option)]}
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
