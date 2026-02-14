<script lang="ts" module>
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLElement>> = $props();
</script>

<span bind:this={ref} data-slot="form-separator" class="select-none" {...restProps}>
	<div class="flex w-full items-center bg-transparent px-2 text-sm">
		{#if children}
			{@render Separator()}
			<span class="px-1 text-center whitespace-pre-wrap text-muted-foreground/40">
				{@render children()}
			</span>
			{@render Separator()}
		{:else}
			{@render Separator()}
		{/if}
	</div>
</span>

{#snippet Separator()}
	<div class="h-px flex-1 bg-muted-foreground/20"></div>
{/snippet}
