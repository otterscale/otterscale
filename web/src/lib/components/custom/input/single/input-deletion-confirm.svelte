<script lang="ts" module>
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import { z } from 'zod';
	import General from './input-general.svelte';
</script>

<script lang="ts">
	type Props = WithElementRef<Omit<HTMLInputAttributes, 'type'>>;

	let {
		ref = $bindable(null),
		value = $bindable(),
		target,
		required,
		class: className,
		...restProps
	}: Props & { target: string } = $props();
</script>

<General
	bind:ref
	required
	type="text"
	schema={z.string().refine(
		(value) => value === target,
		() => ({
			code: 'Unmatch',
			message: `Please type "${target}" to confirm deletion`
		})
	)}
	id="filesystem-delete"
	placeholder={target}
	bind:value
	{...restProps}
/>
