<script lang="ts">
	import { setContext } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Alert from '$lib/components/ui/alert/index';
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';

	import type { AlertType, VariantGetterType } from './types';
	import { IterationManager } from './utils.svelte';
	import { alertVariants } from '../single/alert.svelte';

	let {
		ref = $bindable(null),
		class: className,
		onmouseenter,
		onmouseleave,
		children,
		alerts,
		index = $bindable(0),
		duration = 1000,
		variantGetter,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		alerts: AlertType[];
		index: number;
		duration?: number;
		variantGetter: VariantGetterType;
	} = $props();

	let value = {
		get index() {
			return index;
		},
		set index(newIndex: number) {
			index = newIndex;
		}
	};

	export const iterationManager = new IterationManager(alerts, duration, value, variantGetter);
	setContext('IterationManager', iterationManager);
</script>

<Alert.Root
	bind:ref
	data-slot="alert-root"
	class={cn(
		alertVariants({
			variant: iterationManager.variantGetter(
				iterationManager.alerts[iterationManager.value.index].level
			)
		}),
		'[&>[data-slot=alert-controller]~*]:pr-7',
		className
	)}
	onmouseenter={(e) => {
		iterationManager.stop();
		onmouseenter?.(e);
	}}
	onmouseleave={(e) => {
		iterationManager.start();
		onmouseleave?.(e);
	}}
	{...restProps}
>
	{@render children?.()}
</Alert.Root>
