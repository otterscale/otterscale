<script lang="ts" module>
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import { z } from 'zod';
	import General from './input-general.svelte';

	type Props = WithElementRef<Omit<HTMLInputAttributes, 'type'>>;
</script>

<script lang="ts">
	let {
		id,
		ref = $bindable(null),
		value = $bindable(),
		target,
		required,
		class: className,
		...restProps
	}: Props & { target: string } = $props();

	const schema = z.string().refine(
		(value) => value === target,
		() => ({
			code: 'Unmatch',
			message: `Please type "${target}" to confirm deletion`
		})
	);
</script>

<General
	{id}
	bind:ref
	data-slot="input-delete-confirm"
	required
	type="text"
	{schema}
	placeholder={target}
	bind:value
	{...restProps}
/>
