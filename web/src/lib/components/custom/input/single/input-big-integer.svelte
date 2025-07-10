<script lang="ts" module>
	import InputGeneral from '$lib/components/custom/input/single/input-general.svelte';
	import { cn } from '$lib/utils.js';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { ZodFirstPartySchemaTypes } from 'zod';
</script>

<script lang="ts">
	let {
		value = $bindable(),
		required,
		schema,
		class: className,
		...restProps
	}: Omit<HTMLInputAttributes, 'type'> & { schema?: ZodFirstPartySchemaTypes } = $props();

	const INPUT_TYPE = 'number';
</script>

<InputGeneral
	data-slot="input-big-integer"
	bind:value
	type={INPUT_TYPE}
	{required}
	{schema}
	class={cn(className)}
	oninput={() => {
		value = String(value);
	}}
	{...restProps}
/>
