<script lang="ts">
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';

	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<AlertDialog.Footer
	bind:ref
	data-slot="alert-dialog-footer"
	class={cn('w-full', className)}
	{...restProps}
>
	<div class="flex w-full items-center justify-between">
		{@render children?.()}
	</div>
</AlertDialog.Footer>
