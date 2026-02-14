<script lang="ts" module>
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLInputAttributes } from 'svelte/elements';

	import { cn } from '$lib/utils';

	import General from './input-general.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		value = $bindable(),
		class: className,
		required,
		target,
		invalid = $bindable(),
		...restProps
	}: WithElementRef<Omit<HTMLInputAttributes, 'type' | 'files'>> & {
		target: string;
		invalid?: boolean | null | undefined;
	} = $props();

	const isInvalid = $derived(required && value !== target);
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
</div>
