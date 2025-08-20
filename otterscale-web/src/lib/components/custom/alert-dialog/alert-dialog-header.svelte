<script lang="ts">
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index';
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<AlertDialog.Header
	bind:ref
	class={cn('flex items-center justify-center text-xl font-bold', className)}
	{...restProps}
>
	{@render children?.()}
</AlertDialog.Header>
