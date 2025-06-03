<script lang="ts">
	import { setContext } from 'svelte';
	import { InputManager, ValuesManager } from './utils.svelte';
	import type { InputType } from './types';

	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		values = $bindable(),
		type,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		values: any[];
		type: InputType;
	} = $props();

	setContext('InputManager', new InputManager(type));
	setContext(
		'ValuesManager',
		new ValuesManager(values, (newValues: any[]) => {
			values = newValues;
		})
	);
</script>

<div
	bind:this={ref}
	data-slot="input-controller"
	class={cn('flex flex-col gap-2', className)}
	{...restProps}
>
	{@render children?.()}
</div>
