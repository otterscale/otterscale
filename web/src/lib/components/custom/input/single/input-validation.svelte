<script lang="ts" module>
	import { Badge } from '$lib/components/ui/badge';
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';
	import { z } from 'zod';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		isInvalid,
		errors,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		isInvalid: boolean;
		errors?: z.ZodIssue[];
	} = $props();
</script>

{#if isInvalid && errors}
	<div
		bind:this={ref}
		data-slot="input-required"
		class={cn('animate-in fade-in text-destructive flex items-center gap-2', className)}
		{...restProps}
	>
		{#each errors as error}
			<span class="flex items-center gap-1">
				<Badge variant="destructive">{error.code}</Badge>
				<p class="text-xs">{error.message}</p>
			</span>
		{/each}
	</div>
{/if}
