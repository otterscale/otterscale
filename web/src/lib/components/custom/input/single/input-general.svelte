<script lang="ts">
	import { BORDER_INPUT_CLASSNAME, UNFOCUS_INPUT_CLASSNAME, typeToIcon } from './utils.svelte';
	import InputRequired from './input-required.svelte';
	import InputValidation from './input-validation.svelte';

	import Icon from '@iconify/svelte';
	import { Input } from '$lib/components/ui/input';
	import type { ZodFirstPartySchemaTypes } from 'zod';
	import { InputValidator } from './utils.svelte';
	import type { InputType } from './types';

	import type { HTMLInputAttributes } from 'svelte/elements';
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';

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
				<Icon icon={typeToIcon[type]} />
			</span>
		{/if}
		<Input
			bind:ref
			data-slot="input-general"
			class={cn(UNFOCUS_INPUT_CLASSNAME)}
			{type}
			bind:value
			{...restProps}
		/>
	</div>
{/snippet}

{#if isNotFilled}
	<InputRequired {isNotFilled} />
{/if}
{#if schema}
	{@const validator = new InputValidator(schema)}
	{@const validation = validator.validate(value)}

	{@const isInvalid = value && !validation.isValid}

	{@render Controller([
		BORDER_INPUT_CLASSNAME,
		isNotFilled || isInvalid ? 'ring-destructive ring-1' : ''
	])}
	{#if isInvalid}
		<InputValidation {isInvalid} errors={validation.errors} />
	{/if}
{:else}
	{@render Controller([BORDER_INPUT_CLASSNAME])}
{/if}
