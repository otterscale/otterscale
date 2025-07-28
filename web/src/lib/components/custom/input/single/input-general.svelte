<script lang="ts" module>
	import { FormValidator } from '$lib/components/custom/form';
	import { Input } from '$lib/components/ui/input';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import { getContext, hasContext } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { ZodFirstPartySchemaTypes } from 'zod';
	import InputValidation from './input-validation.svelte';
	import type { InputType } from './types';
	import { InputValidator, typeToIcon } from './utils.svelte';

	type Props = WithElementRef<Omit<HTMLInputAttributes, 'type'> & { type?: InputType }>;
</script>

<script lang="ts">
	let {
		id,
		ref = $bindable(null),
		value = $bindable(),
		type,
		required,
		schema,
		class: className,
		oninput,
		transformer = (value) => value,
		...restProps
	}: Props & {
		schema?: ZodFirstPartySchemaTypes;
		transformer?: (value: any) => void;
	} = $props();

	const isNotFilled = $derived(required && (value === '' || value === undefined));

	const validator = new InputValidator(schema);
	const validation = $derived(validator.validate(value));
	const isInvalid = $derived(value && !validation.isValid);

	const formValidator: FormValidator = getContext('FormValidator');
</script>

<div class="relative">
	{#if type}
		<span class="absolute left-3 top-1/2 -translate-y-1/2 items-center">
			<Icon icon={hasContext('icon') ? getContext('icon') : typeToIcon[type]} />
		</span>
	{/if}

	<Input
		bind:ref
		data-slot="input-general"
		placeholder={isNotFilled ? 'Required' : ''}
		{type}
		bind:value
		oninput={(e) => {
			value = transformer(value);
			formValidator.set(id, isNotFilled || isInvalid);
			oninput?.(e);
		}}
		class={cn(
			'pl-9 ring-1',
			isNotFilled ? 'placeholder:text-destructive/60 placeholder:text-xs' : '',
			isNotFilled || isInvalid ? 'ring-destructive' : '',
			className
		)}
		{...restProps}
	/>
</div>

{#if isInvalid}
	<InputValidation {isInvalid} errors={validation.errors} />
{/if}
