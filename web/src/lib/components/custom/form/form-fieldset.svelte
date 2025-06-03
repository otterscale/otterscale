<script lang="ts">
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLFieldsetAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLFieldsetAttributes> = $props();
</script>

<fieldset
	bind:this={ref}
	data-slot="form-fieldset"
	class={cn('space-y-4 rounded-md border px-4 pb-4', className)}
	{...restProps}
>
	{@render children?.()}
</fieldset>
