<script lang="ts">
	import * as Alert from '$lib/components/ui/alert';
	import type { HTMLAttributes } from 'svelte/elements';
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<Alert.Title bind:ref data-slot="alert-title" class={cn('text-lg', className)} {...restProps}>
	{@render children?.()}
</Alert.Title>
