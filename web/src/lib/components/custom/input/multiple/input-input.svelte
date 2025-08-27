<script lang="ts" module>
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';
	import { getContext } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import * as Input from '../single';
	import { InputManager, ValuesManager } from './utils.svelte';
	import { validate } from './utils.svelte';

	type Props = WithElementRef<Omit<HTMLInputAttributes, 'value' | 'type'>>;
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		onkeydown,
		onblur,
		...restProps
	}: Props = $props();

	const id: string | undefined = getContext('id');
	const required: boolean | undefined = getContext('required');
	const inputManager: InputManager = getContext('InputManager');
	const valuesManager: ValuesManager = getContext('ValuesManager');

	const isInvalid = $derived(validate(required, valuesManager));
</script>

<div class="w-full">
	<Input.General
		{id}
		bind:ref
		data-slot="input-input"
		type={inputManager.type}
		bind:value={inputManager.input}
		class={cn(
			'ring-1',
			isInvalid
				? 'placeholder:text-destructive/60 placeholder:text-xs focus:placeholder:invisible'
				: '',
			isInvalid ? 'ring-destructive' : '',
			className
		)}
		placeholder={isInvalid ? 'Required' : ''}
		onkeydown={(e) => {
			if (e.key === 'Enter') {
				valuesManager.append(inputManager.input);
				inputManager.reset();
				onkeydown?.(e);
			}
		}}
		{...restProps}
	/>
</div>
