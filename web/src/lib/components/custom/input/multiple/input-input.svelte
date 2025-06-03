<script lang="ts">
	import { getContext } from 'svelte';
	import * as Input from '../single';
	import { InputManager, ValuesManager } from './utils.svelte';
	import type { InputType } from './types';

	import type { HTMLInputAttributes } from 'svelte/elements';
	import { cn, type WithElementRef } from '$lib/utils.js';

	type Props = WithElementRef<Omit<HTMLInputAttributes, 'value' | 'type'> & { type?: InputType }>;

	let {
		ref = $bindable(null),
		class: className,
		onkeydown,
		onblur,
		...restProps
	}: Props & {} = $props();

	const { files, ...restProperties } = restProps;
	const inputManager: InputManager = getContext('InputManager');
	const valuesManager: ValuesManager = getContext('ValuesManager');
</script>

<div class="w-full">
	{#if inputManager.type === 'color'}
		<Input.Color
			bind:ref
			data-slot="input-input"
			type={inputManager.type}
			bind:value={inputManager.input}
			class={cn('w-full', className)}
			onkeydown={(e) => {
				if (e.key === 'Enter') {
					valuesManager.append(e.currentTarget.value);
					inputManager.reset();
					onkeydown?.(e);
				}
			}}
			{...restProperties}
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
			{...restProperties}
		/>
	{/if}
</div>
