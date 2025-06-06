<script lang="ts">
	import { setContext } from 'svelte';
	import { InputManager, ValuesManager } from './utils.svelte';
	import type { InputType } from './types';

	import { cn } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';
	import type { WithElementRef } from 'bits-ui';
	import type { string } from 'zod';

	let {
		ref = $bindable(null),
		class: className,
		values = $bindable(),
		type,
		contextData,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		values: any[];
		type: InputType;
		contextData?: Record<string, any>;
	} = $props();

	setContext('InputManager', new InputManager(type));
	setContext(
		'ValuesManager',
		new ValuesManager(values, (newValues: any[]) => {
			values = newValues;
		})
	);
	if (contextData) {
		for (const key in contextData) {
			setContext(key, contextData[key]);
		}
	}
</script>

<div
	bind:this={ref}
	data-slot="input-controller"
	class={cn('flex flex-col gap-2', className)}
	{...restProps}
>
	{@render children?.()}
</div>
