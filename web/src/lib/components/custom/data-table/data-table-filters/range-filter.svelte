<script lang="ts" module>
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { Slider } from '$lib/components/ui/slider/index.js';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';
	import { capitalizeFirstLetter } from 'better-auth';
</script>

<script lang="ts" generics="TData">
	let { table, columnId, alias }: { table: Table<TData>; columnId: string; alias?: string } =
		$props();

	const values =
		table.getCoreRowModel().rows.map((row) => Number(row.getValue(columnId))) ?? ([] as number[]);

	const min = Math.min(...values);
	const max = Math.max(...values);

	let rangeMin = $state(min);
	let rangeMax = $state(max);

	let value = $derived([rangeMin, rangeMax]);
</script>

<Popover.Root>
	<Popover.Trigger>
		<Button variant="outline"
			><Icon icon="ph:funnel" />{alias ?? capitalizeFirstLetter(columnId)}</Button
		>
	</Popover.Trigger>
	<Popover.Content
		class={cn(
			'flex flex-col items-center justify-between gap-4 sm:w-[62vw] md:w-[50vw] lg:w-[38vw]'
		)}
	>
		<Slider
			type="multiple"
			{value}
			{min}
			{max}
			step={1}
			onValueChange={(v) => {
				rangeMin = v[0];
				rangeMax = v[1];
				table.getColumn(columnId)?.setFilterValue(value);
			}}
		/>
		<div class="flex w-full items-center justify-between gap-2">
			<Input
				class="text-xs"
				bind:value={rangeMin}
				type="number"
				oninput={() => {
					table.getColumn(columnId)?.setFilterValue(value);
				}}
			/>
			<p class="text-sm">to</p>
			<Input
				class="text-xs"
				bind:value={rangeMax}
				type="number"
				oninput={() => {
					table.getColumn(columnId)?.setFilterValue(value);
				}}
			/>
		</div>
	</Popover.Content>
</Popover.Root>
