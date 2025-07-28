<script lang="ts" module>
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLFieldsetAttributes } from 'svelte/elements';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		disabled,
		children,
		...restProps
	}: WithElementRef<HTMLFieldsetAttributes> = $props();
</script>

<fieldset
	bind:this={ref}
	data-slot="form-fieldset"
	class={cn(
		'select-none space-y-4 rounded-md border p-4',
		disabled ? 'bg-muted/40' : '',
		className
	)}
	{...restProps}
>
	{@render children?.()}
</fieldset>
