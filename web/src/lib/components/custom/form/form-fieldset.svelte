<script lang="ts">
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';
	import Icon from '@iconify/svelte';
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
	class={cn(
		'space-y-4 rounded-md border px-4 pb-4',
		restProps.disabled ? 'bg-muted/40' : '',
		className
	)}
	{...restProps}
>
	{#if restProps.disabled}
		<div class="text-destructive flex items-center justify-center gap-1">
			<Icon icon="ph:warning" />
			<p>Disabled</p>
		</div>
	{/if}
	{@render children?.()}
</fieldset>
