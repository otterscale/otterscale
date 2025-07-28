<script lang="ts" module>
	import * as Select from '$lib/components/ui/select';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { ZodFirstPartySchemaTypes } from 'zod';
	import { General } from '.';
	import type { InputType, UnitType } from './types';

	type Props = WithElementRef<Omit<HTMLInputAttributes, 'type'> & { type?: InputType }>;
</script>

<script lang="ts">
	let {
		id,
		ref = $bindable(null),
		value = $bindable(),
		required,
		schema,
		units,
		transformer = (value) => value,
		oninput,
		class: className,
		...restProps
	}: Props & {
		schema?: ZodFirstPartySchemaTypes;
		units: UnitType[];
		transformer?: (value: any) => void;
	} = $props();

	function getDefault(): { value: number | undefined; unit: UnitType | undefined } {
		const UNITS = units.sort((p, n) => p.value - n.value);

		const INITIAL_VALUE = value !== undefined ? Number(value) : undefined;

		if (INITIAL_VALUE === undefined) {
			return { value: undefined, unit: UNITS[0] };
		}

		let temporaryValue = 0;
		let [temporaryUnit] = units;
		for (const unit of UNITS) {
			if (INITIAL_VALUE / unit.value >= 1) {
				temporaryValue = INITIAL_VALUE / unit.value;
				temporaryUnit = unit;
			}
		}
		return { value: temporaryValue, unit: temporaryUnit };
	}

	const DEFAULT = getDefault();
	const DEFAULT_VALUE = DEFAULT.value;
	const DEFAULT_UNIT = DEFAULT.unit;

	let inputValue: number | undefined = $state(DEFAULT_VALUE);
	let unit: UnitType | undefined = $state(DEFAULT_UNIT);
</script>

<div class="flex items-center gap-2">
	<div class={cn('w-full')}>
		<General
			bind:ref
			data-slot="input-general"
			type="number"
			bind:value={inputValue}
			{required}
			{...restProps}
			oninput={(e) => {
				value = transformer(inputValue !== undefined && unit ? inputValue * unit.value : undefined);
				oninput?.(e);
			}}
		/>
	</div>
	<Select.Root type="single">
		<Select.Trigger class="w-fit">
			{unit && unit.label ? unit.label : 'No Unit'}
		</Select.Trigger>
		<Select.Content>
			{#each units as option}
				<Select.Item
					value={option.value}
					class="flex items-center gap-2 text-xs hover:cursor-pointer"
					onclick={() => {
						unit = option;
						value = transformer(
							inputValue !== undefined && unit ? inputValue * unit.value : undefined
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
