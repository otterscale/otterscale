<script lang="ts" module>
	import { Input } from '$lib/components/ui/input';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import { getContext, hasContext } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { ZodFirstPartySchemaTypes } from 'zod';
</script>

<script lang="ts">
	import InputValidation from './input-validation.svelte';
	import type { InputType } from './types';
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
		type,
		required,
		schema,
		class: className,
		...restProps
	}: Props & { schema?: ZodFirstPartySchemaTypes } = $props();

	const isNotFilled = $derived(required && !value);
</script>

{#snippet Controller(classNames: string[])}
	<div class={cn(...classNames, className)}>
		{#if type}
			<span class="pl-3">
				<Icon icon={hasContext('icon') ? getContext('icon') : typeToIcon[type]} />
			</span>
		{/if}

		<Input
			bind:ref
			data-slot="input-general"
			placeholder={isNotFilled ? 'Required' : ''}
			class={cn(
				UNFOCUS_INPUT_CLASSNAME,
				isNotFilled ? 'placeholder:text-destructive/60 placeholder:text-xs' : ''
			)}
			{type}
			bind:value
			{...restProps}
		/>
	</div>
{/snippet}

{#if schema}
	{@const validator = new InputValidator(schema)}
	{@const validation = validator.validate(value)}

	{@const isInvalid = value && !validation.isValid}

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
