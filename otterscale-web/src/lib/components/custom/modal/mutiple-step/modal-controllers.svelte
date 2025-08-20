<script lang="ts">
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
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

<AlertDialog.ActionsGroup
	bind:ref
	data-slot="multiple-step-modal-action-group"
	class={cn(className)}
	{...restProps}
>
	{@render children?.()}
</AlertDialog.ActionsGroup>
