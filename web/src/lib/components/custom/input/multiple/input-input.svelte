<script lang="ts" module>
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';
	import { getContext } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
</script>

<script lang="ts">
	import * as Input from '../single';
	import { InputManager, ValuesManager } from './utils.svelte';

	type Props = WithElementRef<Omit<HTMLInputAttributes, 'value' | 'type'>>;

	let {
		ref = $bindable(null),
		class: className,
		onkeydown,
		onblur,
		...restProps
	}: Props = $props();

	const inputManager: InputManager = getContext('InputManager');
	const valuesManager: ValuesManager = getContext('ValuesManager');
	const id: string | undefined = getContext('id');
	const required: boolean | undefined = getContext('required');

	const isNotFilled = $derived(required && valuesManager.values.length === 0);
</script>

<div class="w-full">
	{#if inputManager.type === 'color'}
		<Input.Color
			{id}
			{required}
			bind:ref
			data-slot="input-input"
			bind:value={inputManager.input}
			class={cn('w-full', className)}
			onkeydown={(e) => {
				if (e.key === 'Enter') {
					valuesManager.append(e.currentTarget.value);
					inputManager.reset();
					onkeydown?.(e);
				}
			}}
			{...restProps}
		/>
	{:else}
		<Input.General
			{id}
			bind:ref
			data-slot="input-input"
			type={inputManager.type}
			bind:value={inputManager.input}
			class={cn('ring-1', isNotFilled ? 'ring-destructive' : '', className)}
			onkeydown={(e) => {
				if (e.key === 'Enter') {
					valuesManager.append(inputManager.input);
					inputManager.reset();
					onkeydown?.(e);
				}
			}}
			{...restProps}
		/>
	{/if}
</div>
