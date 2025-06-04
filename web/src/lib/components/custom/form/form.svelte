<script lang="ts">
	import { cn } from '$lib/utils.js';
	import type { HTMLFormAttributes } from 'svelte/elements';
	import type { WithElementRef } from 'bits-ui';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLFormAttributes> = $props();
</script>

<form
	bind:this={ref}
	data-slot="form-root"
	class={cn('flex h-[77vh] flex-col gap-4 overflow-auto rounded-md', className)}
	{...restProps}
>
	{@render children?.()}
</form>
