<script lang="ts">
	import { getContext } from 'svelte';
	import * as Input from '../single';
	import { InputManager, ValuesManager } from './utils.svelte';
	import type { InputType } from './types';

	import type { HTMLInputAttributes } from 'svelte/elements';
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';

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
</script>

<div class="w-full">
	{#if inputManager.type === 'color'}
		<Input.Color
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
			bind:ref
			data-slot="input-input"
			type={inputManager.type}
			bind:value={inputManager.input}
			class={cn(className)}
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
