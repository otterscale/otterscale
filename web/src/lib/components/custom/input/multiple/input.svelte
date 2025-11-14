<script lang="ts" module>
	import type { WithElementRef } from 'bits-ui';
	import { setContext } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	import { cn } from '$lib/utils.js';

	import type { InputType } from './types';
	import { InputManager, validate, ValuesManager } from './utils.svelte';
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
		icon,
		invalid = $bindable(),
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		values?: any[];
		required?: boolean;
		type?: InputType;
		icon?: string;
		invalid?: boolean | null | undefined;
	} = $props();
	const inputManager = new InputManager(type, icon);
	const valuesManager = new ValuesManager(values, {
		set values(newValue: any[]) {
			values = newValue;
		}
	});
	setContext('id', id);
	setContext('required', required);
	setContext('InputManager', inputManager);
	setContext('ValuesManager', valuesManager);

	$effect(() => {
		invalid = validate(required, valuesManager);
	});
</script>

<div
	bind:this={ref}
	data-slot="input-controller"
	class={cn('flex flex-col gap-2', className)}
	{...restProps}
>
	{@render children?.()}
</div>
