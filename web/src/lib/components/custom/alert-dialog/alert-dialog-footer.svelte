<script lang="ts" module>
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<AlertDialog.Footer bind:ref data-slot="alert-dialog-footer" class={cn('w-full', className)} {...restProps}>
	<div class="flex w-full items-center justify-between gap-2">
		{@render children?.()}
	</div>
</AlertDialog.Footer>
