<script lang="ts">
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
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

<AlertDialog.Footer
	bind:ref
	data-slot="alert-dialog-actions-group"
	class={cn('w-full', className)}
	{...restProps}
>
	<div class="flex items-center justify-end gap-2">
		{@render children?.()}
	</div>
</AlertDialog.Footer>
