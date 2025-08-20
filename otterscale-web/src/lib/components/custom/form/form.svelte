<script lang="ts">
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLFormAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		children,
		invalid = $bindable(),
		...restProps
	}: WithElementRef<HTMLFormAttributes> & { invalid?: boolean | undefined } = $props();
</script>

<form
	bind:this={ref}
	data-slot="form-root"
	class={cn('flex max-h-[77vh] flex-col gap-4 overflow-auto rounded-md', className)}
	{...restProps}
>
	{@render children?.()}
</form>
