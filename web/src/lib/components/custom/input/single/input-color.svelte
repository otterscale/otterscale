<script lang="ts" module>
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { Input } from '$lib/components/ui/input';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import { z, type ZodFirstPartySchemaTypes } from 'zod';
</script>

<script lang="ts">
	import InputValidation from './input-validation.svelte';
	import {
		BORDER_INPUT_CLASSNAME,
		InputValidator,
		RING_INVALID_INPUT_CLASSNAME,
		RING_VALID_INPUT_CLASSNAME,
		typeToIcon,
		UNFOCUS_INPUT_CLASSNAME
	} from './utils.svelte';

	let {
		ref = $bindable(null),
		value = $bindable(),
		schema = z.string().regex(/^#[0-9a-fA-F]{6}$/),
		class: className,
		...restProps
	}: WithElementRef<Exclude<HTMLInputAttributes, 'type'>> & { type?: 'color' } & {
		schema?: ZodFirstPartySchemaTypes;
	} = $props();

	const validator = new InputValidator(schema);
	const validation = $derived(validator.validate(value));
	const isInvalid = $derived(value && !validation.isValid);
</script>

<div
	class={cn(
		BORDER_INPUT_CLASSNAME,
		isInvalid ? RING_INVALID_INPUT_CLASSNAME : RING_VALID_INPUT_CLASSNAME,
		'h-10 justify-between',
		className
	)}
>
	<span class="flex items-center gap-2">
		<span class="pl-3">
			<Icon icon={typeToIcon['color']} />
		</span>
		<Badge variant="outline">{value}</Badge>
	</span>

	<Input
		bind:ref
		data-slot="input-color"
		class={cn(UNFOCUS_INPUT_CLASSNAME, 'mr-3 aspect-square h-7 w-fit cursor-pointer p-0')}
		type="color"
		bind:value
		{...restProps}
	/>
</div>

<div class="transition-all duration-500">
	<InputValidation {isInvalid} errors={validation.errors} />
</div>
