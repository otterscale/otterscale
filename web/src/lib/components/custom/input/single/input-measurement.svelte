<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import { onMount } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';

	import * as Select from '$lib/components/ui/select';
	import { cn } from '$lib/utils.js';

	import { General } from '.';
	import type { InputType, UnitType } from './types';
	import { getMeasurement } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		value = $bindable(),
		required,
		units,
		oninput,
		transformer = (value) => value,
		invalid = $bindable(),
		...restProps
	}: WithElementRef<Omit<HTMLInputAttributes, 'type'> & { type?: InputType }> & {
		units: UnitType[];
		transformer?: (value: any) => void;
		invalid?: boolean | null | undefined;
	} = $props();

	let temporaryValue: number | undefined = $state(undefined);
	let temporaryUnit: UnitType | undefined = $state(undefined);
	onMount(() => {
		const measurement = getMeasurement(value, units);
		temporaryValue = measurement.value;
		temporaryUnit = measurement.unit;
	});

	const isInvalid = $derived(required && (value === null || value === undefined));
	$effect(() => {
		invalid = isInvalid;
	});
</script>

<div class="flex items-center gap-2">
	<div class={cn('w-full')}>
		<General
			bind:ref
			type="number"
			bind:value={temporaryValue}
			{required}
			{invalid}
			oninput={(e) => {
				value = transformer(
					typeof temporaryValue === 'number' && temporaryUnit !== undefined
						? temporaryValue * temporaryUnit.value
						: undefined
				);
				oninput?.(e);
			}}
			{...restProps}
		/>
	</div>
	<Select.Root type="single" value={temporaryUnit?.value}>
		<Select.Trigger class="w-fit">
			{temporaryUnit && temporaryUnit.label ? temporaryUnit.label : 'No Unit'}
		</Select.Trigger>
		<Select.Content>
			{#each units as option (option.value)}
				<Select.Item
					value={option.value}
					class="flex items-center gap-2 text-xs hover:cursor-pointer"
					onclick={() => {
						temporaryUnit = option;
						value = transformer(
							typeof temporaryValue === 'number' && temporaryUnit !== undefined
								? temporaryValue * temporaryUnit.value
								: undefined
						);
					}}
				>
					<Icon icon={option.icon ?? 'ph:scales'} class={cn('size-4')} />
					{option.label}
				</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>
</div>
