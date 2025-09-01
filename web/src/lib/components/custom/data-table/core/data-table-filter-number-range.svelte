<script lang="ts" module>
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { Slider } from '$lib/components/ui/slider/index.js';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	const empty = [undefined, undefined];
</script>

<script lang="ts" generics="TData">
	let {
		table,
		values,
		columnId,
		alias,
	}: { table: Table<TData>; values: number[]; columnId: string; alias?: string } = $props();

	const minimum = Math.min(...values);
	const maximum = Math.max(...values);

	let rangeMinimum: number | undefined = $state(minimum);
	let rangeMaximum: number | undefined = $state(maximum);

	let value = $derived([rangeMinimum, rangeMaximum]);
</script>

<Popover.Root>
	<Popover.Trigger class={cn(buttonVariants({ size: 'sm', variant: 'outline' }), 'text-xs')}>
		{@render filterTrigger()}
	</Popover.Trigger>
	<Popover.Content class="flex w-[500px] flex-col items-center justify-center gap-4">
		<div class="flex w-full items-center justify-between gap-2">
			<Input
				class="h-7 w-16 text-xs"
				bind:value={rangeMinimum}
				type="number"
				oninput={() => {
					table.getColumn(columnId)?.setFilterValue(value);
				}}
			/>
			<Slider
				type="multiple"
				{value}
				min={minimum}
				max={maximum}
				step={1}
				onValueChange={(v) => {
					[rangeMinimum, rangeMaximum] = v;
					table.getColumn(columnId)?.setFilterValue(value);
				}}
			/>
			<Input
				class="h-7 w-16 text-xs"
				bind:value={rangeMaximum}
				type="number"
				oninput={() => {
					table.getColumn(columnId)?.setFilterValue(value);
				}}
			/>
		</div>
		{@render filterActions()}
	</Popover.Content>
</Popover.Root>

{#snippet filterTrigger()}
	<Icon icon="ph:funnel" />
	<span class="capitalize">{alias ?? columnId}</span>
	<Separator orientation="vertical" />
	<Badge variant="outline" class="size-sm">
		{rangeMinimum}~{rangeMaximum}
	</Badge>
{/snippet}

{#snippet filterActions()}
	<div class="items-cent1er flex w-full items-center justify-end gap-2">
		<Button
			variant="outline"
			size="sm"
			onclick={() => {
				rangeMinimum = minimum;
				rangeMaximum = maximum;
				table.getColumn(columnId)?.setFilterValue(value);
			}}
		>
			{m.all()}
		</Button>
		<Button
			variant="secondary"
			size="sm"
			onclick={() => {
				rangeMinimum = minimum;
				rangeMaximum = minimum;
				table.getColumn(columnId)?.setFilterValue(empty);
			}}
		>
			{m.clear()}
		</Button>
	</div>
{/snippet}
