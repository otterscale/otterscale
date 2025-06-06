<script lang="ts" generics="TData">
	import { cn } from '$lib/utils.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import Icon from '@iconify/svelte';
	import { Input } from '$lib/components/ui/input';

	import { type Table } from '@tanstack/table-core';
	import Separator from '$lib/components/ui/separator/separator.svelte';

	import { Slider } from '$lib/components/ui/slider/index.js';

	let { table }: { table: Table<TData> } = $props();

	const amounts = table.getCoreRowModel().rows.map((row) => row.getValue('amount')) as number[];
	const min = Math.min(...amounts);
	const max = Math.max(...amounts);
	const Range = [min, max];
	const Q1 = min + (max - min) / 4;
	const Q3 = min + (3 * (max - min)) / 4;
	const IQR = [Q1, Q3];

	let rangeMin = $state(min);
	let rangeMax = $state(max);
	let value = $derived([rangeMin, rangeMax]);
</script>

<Popover.Root>
	<Popover.Trigger>
		<Button variant="outline"><Icon icon="ph:funnel" />Amount</Button>
	</Popover.Trigger>
	<Popover.Content
		class={cn(
			'flex flex-col items-center justify-between gap-4 sm:w-[62vw] md:w-[50vw] lg:w-[38vw]'
		)}
	>
		<div class="flex w-full flex-wrap gap-2">
			<div class="flex items-center gap-2">
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						<Badge variant="outline" class="hover:bg-muted h-6 text-xs hover:cursor-pointer">
							Range
						</Badge>
					</DropdownMenu.Trigger>
					<DropdownMenu.Content>
						<DropdownMenu.Item
							class="text-xs hover:cursor-pointer"
							onclick={() => {
								rangeMin = min;
								rangeMax = max;
								table.getColumn('amount')?.setFilterValue(Range);
							}}
						>
							ALL
						</DropdownMenu.Item>
						<DropdownMenu.Item
							class="text-xs hover:cursor-pointer"
							onclick={() => {
								rangeMin = Q1;
								rangeMax = Q3;
								table.getColumn('amount')?.setFilterValue(IQR);
							}}
						>
							IQR
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						<Badge variant="outline" class="hover:bg-muted h-6 text-xs hover:cursor-pointer">
							Top
						</Badge>
					</DropdownMenu.Trigger>
					<DropdownMenu.Content>
						{#each [10, 20, 30, 40, 50, 60, 70, 80, 90, 100] as n}
							<DropdownMenu.Item
								class="text-xs hover:cursor-pointer"
								onclick={() => {
									const sorted = [...amounts].sort((a, b) => b - a);
									const threshold =
										sorted[Math.min(n - 1, Math.round(sorted.length * (n / 100)) - 1)];
									rangeMin = threshold;
									rangeMax = max;
									table.getColumn('amount')?.setFilterValue([threshold, max]);
								}}
							>
								{n}%
							</DropdownMenu.Item>
						{/each}
					</DropdownMenu.Content>
				</DropdownMenu.Root>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						<Badge variant="outline" class="hover:bg-muted h-6 text-xs hover:cursor-pointer">
							Tail
						</Badge>
					</DropdownMenu.Trigger>
					<DropdownMenu.Content>
						{#each [10, 20, 30, 40, 50, 60, 70, 80, 90, 100] as n}
							<DropdownMenu.Item
								class="text-xs hover:cursor-pointer"
								onclick={() => {
									const sorted = [...amounts].sort((a, b) => a - b);
									const threshold =
										sorted[Math.min(n - 1, Math.round(sorted.length * (n / 100)) - 1)];
									rangeMin = min;
									rangeMax = threshold;
									table.getColumn('amount')?.setFilterValue([min, threshold]);
								}}
							>
								{n}%
							</DropdownMenu.Item>
						{/each}
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</div>
		</div>
		<div class="flex w-full flex-col gap-2">
			<div class="text-muted-foreground flex w-full justify-between text-xs font-light">
				{#each Array.from({ length: 5 }, (_, i) => i) as index}
					{@const interval = (max - min) / 4}
					<div class="flex flex-col items-center justify-center gap-1">
						{Math.round(min + index * interval)}
					</div>
				{/each}
			</div>
			<Slider
				type="multiple"
				{value}
				{min}
				{max}
				step={1}
				onValueChange={(v) => {
					rangeMin = v[0];
					rangeMax = v[1];
					table.getColumn('amount')?.setFilterValue(value);
				}}
			/>
		</div>
		<div class="flex w-full items-center justify-between gap-2">
			<Input
				class="text-xs"
				bind:value={rangeMin}
				type="number"
				oninput={() => {
					table.getColumn('amount')?.setFilterValue(value);
				}}
			/>
			<p class="text-sm">to</p>
			<Input
				class="text-xs"
				bind:value={rangeMax}
				type="number"
				oninput={() => {
					table.getColumn('amount')?.setFilterValue(value);
				}}
			/>
		</div>
	</Popover.Content>
</Popover.Root>
