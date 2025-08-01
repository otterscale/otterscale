<script lang="ts" module>
	import { cn } from '$lib/utils';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import General from './input-general.svelte';
	import { Copy } from '$lib/components/custom/copy';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		value = $bindable(),
		class: className,
		id,
		required,
		target,
		invalid = $bindable(),
		...restProps
	}: WithElementRef<Omit<HTMLInputAttributes, 'type' | 'files'>> & {
		target: string;
		invalid?: boolean | null | undefined;
	} = $props();

	const isInvalid = $derived(value !== target);
	$effect(() => {
		invalid = isInvalid;
	});
</script>

<div class="relative">
	<General
		bind:ref
		data-slot="input-delete-confirm"
		class={cn('pr-9', isInvalid ? 'ring-destructive' : '', className)}
		type="text"
		bind:value
		placeholder={target}
		required
		{...restProps}
	/>
	<Copy text={target} class="absolute right-0 top-1/2 -translate-y-1/2 items-center" />
</div>
