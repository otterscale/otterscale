<script lang="ts" module>
	import type { WithElementRef } from 'bits-ui';
	import { setContext } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Alert from '$lib/components/ui/alert/index';
	import { cn } from '$lib/utils.js';

	import { alertVariants } from '../single/alert.svelte';
	import type { AlertType } from './types';
	import { IterationManager } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		onmouseenter,
		onmouseleave,
		children,
		alerts,
		index = $bindable(0),
		duration = 1000,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		alerts: AlertType[];
		index: number;
		duration?: number;
	} = $props();

	let value = {
		get index() {
			return index;
		},
		set index(newIndex: number) {
			index = newIndex;
		}
	};

	export const iterationManager = new IterationManager(alerts, duration, value);
	setContext('IterationManager', iterationManager);
</script>

<Alert.Root
	bind:ref
	data-slot="alert-root"
	class={cn(
		alertVariants({
			variant: alerts[index].variant
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
