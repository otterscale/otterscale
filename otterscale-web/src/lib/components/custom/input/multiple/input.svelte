<script lang="ts" module>
	import { setContext } from 'svelte';
	import { InputManager, ValuesManager } from './utils.svelte';
	import type { InputType } from './types';
	import { cn } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';
	import type { WithElementRef } from 'bits-ui';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		values = $bindable(),
		class: className,
		children,
		id,
		required,
		type,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		values: any[];
		required?: boolean;
		type: InputType;
	} = $props();

	setContext('id', id);
	setContext('required', required);
	setContext('InputManager', new InputManager(type));
	setContext(
		'ValuesManager',
		new ValuesManager(values, {
			set values(newValue: any[]) {
				values = newValue;
			}
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
