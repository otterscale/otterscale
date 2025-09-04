<script lang="ts" module>
	import * as Select from '$lib/components/ui/select';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import { General } from '.';
	import type { InputType, UnitType } from './types';
	import { getInputMeasurementUnitByValue } from './utils.svelte';
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

	const DEFAULT = getInputMeasurementUnitByValue(value, units);
	const DEFAULT_VALUE = DEFAULT.value;
	const DEFAULT_UNIT = DEFAULT.unit;

	let temporaryValue: number | undefined = $state(DEFAULT_VALUE);
	let unit: UnitType | undefined = $state(DEFAULT_UNIT);

	const isInvalid = $derived(required && (value === null || value === undefined));
	$effect(() => {
		invalid = isInvalid;
	});
</script>

<div class="flex items-center gap-2">
	<div class={cn('w-full')}>
		<General
			bind:ref
			data-slot="input-general"
			type="number"
			bind:value={temporaryValue}
			{required}
			oninput={(e) => {
				value = transformer(temporaryValue && unit ? temporaryValue * unit.value : undefined);
				oninput?.(e);
			}}
			{...restProps}
		/>
	</div>
	<Select.Root type="single">
		<Select.Trigger class={cn('w-fit')}>
			{unit && unit.label ? unit.label : 'No Unit'}
		</Select.Trigger>
		<Select.Content>
			{#each units as option}
				<Select.Item
					value={option.value}
					class="flex items-center gap-2 text-xs hover:cursor-pointer"
					onclick={() => {
						unit = option;
						value = transformer(temporaryValue && unit ? temporaryValue * unit.value : undefined);
					}}
				>
					<Icon icon={option.icon ?? 'ph:scales'} class={cn('size-4')} />
					{option.label}
				</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>
</div>
