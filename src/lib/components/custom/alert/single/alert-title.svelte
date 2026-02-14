<script lang="ts">
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Alert from '$lib/components/ui/alert';
	import { cn } from '$lib/utils.js';

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
