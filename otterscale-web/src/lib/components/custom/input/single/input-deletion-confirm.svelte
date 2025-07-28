<script lang="ts" module>
	import { FormValidator } from '$lib/components/custom/form';
	import { cn } from '$lib/utils';
	import type { WithElementRef } from 'bits-ui';
	import { getContext } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import General from './input-general.svelte';

	type Props = WithElementRef<Omit<HTMLInputAttributes, 'type' | 'files'>>;
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		value = $bindable(),
		class: className,
		id,
		required,
		target,
		...restProps
	}: Props & { target: string } = $props();

	const isInvalid = $derived(value !== target);

	const formValidator: FormValidator = getContext('FormValidator');
	$effect(() => {
		formValidator.set(id, isInvalid);
	});
</script>

<General
	bind:ref
	data-slot="input-delete-confirm"
	class={cn(isInvalid ? 'ring-destructive' : '', className)}
	type="text"
	bind:value
	{id}
	placeholder={target}
	required
	{...restProps}
/>
