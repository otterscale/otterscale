<script lang="ts" module>
	import { Input } from '$lib/components/ui/input';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import { getContext, hasContext } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { ZodFirstPartySchemaTypes } from 'zod';
	import * as Select from '$lib/components/ui/select';

	const TYPE = 'number';
</script>

<script lang="ts">
	import InputValidation from './input-validation.svelte';
	import type { InputType, UnitType } from './types';
	import {
		BORDER_INPUT_CLASSNAME,
		InputValidator,
		RING_INVALID_INPUT_CLASSNAME,
		RING_VALID_INPUT_CLASSNAME,
		typeToIcon,
		UNFOCUS_INPUT_CLASSNAME
	} from './utils.svelte';
	type Props = WithElementRef<Omit<HTMLInputAttributes, 'type'> & { type?: InputType }>;

	let {
		ref = $bindable(null),
		value = $bindable(),
		required,
		schema,
		units,
		class: className,
		...restProps
	}: Props & { schema?: ZodFirstPartySchemaTypes; units: UnitType[] } = $props();

	let inputValue = $state(undefined);
	const isNotFilled = $derived(required && !inputValue);

	let unit: UnitType = $state(units.length > 0 ? units[0] : ({} as UnitType));
</script>

{#snippet Controller(classNames: string[])}
	<div class="flex items-center gap-2">
		<div class={cn('w-full', ...classNames, className)}>
			<span class="pl-3">
				<Icon icon={hasContext('icon') ? getContext('icon') : typeToIcon[TYPE]} />
			</span>

			<Input
				bind:ref
				data-slot="input-general"
				placeholder={isNotFilled ? 'Required' : ''}
				class={cn(
					UNFOCUS_INPUT_CLASSNAME,
					isNotFilled ? 'placeholder:text-destructive/60 placeholder:text-xs' : ''
				)}
				type={TYPE}
				bind:value={inputValue}
				{...restProps}
				oninput={() => {
					value = inputValue ? inputValue * unit.value : undefined;
				}}
			/>
		</div>
		<Select.Root type="single">
			<Select.Trigger class="w-fit">{unit && unit.label ? unit.label : 'No Unit'}</Select.Trigger>
			<Select.Content>
				{#each units as option}
					<Select.Item
						value={option.value}
						class="hover:cursor-pointer"
						onclick={() => {
							unit = option;
						}}
					>
						{#if option.icon}
							<Icon
								icon={option.icon}
								class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
							/>
						{/if}
						{option.label}
					</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
	</div>
{/snippet}

{#if schema}
	{@const validator = new InputValidator(schema)}
	{@const validation = validator.validate(inputValue)}

	{@const isInvalid = inputValue && !validation.isValid}

	{@render Controller([
		BORDER_INPUT_CLASSNAME,
		isNotFilled || isInvalid ? RING_INVALID_INPUT_CLASSNAME : RING_VALID_INPUT_CLASSNAME
	])}

	{#if isInvalid}
		<InputValidation {isInvalid} errors={validation.errors} />
	{/if}
{:else}
	{@render Controller([
		BORDER_INPUT_CLASSNAME,
		isNotFilled ? RING_INVALID_INPUT_CLASSNAME : RING_VALID_INPUT_CLASSNAME
	])}
{/if}
