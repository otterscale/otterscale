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

<Alert.Description bind:ref data-slot="alert-description" class={cn('my-2 text-pretty', className)} {...restProps}>
	{@render children?.()}
</Alert.Description>
