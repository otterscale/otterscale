<script lang="ts" module>
	import { FormValidator } from '$lib/components/custom/form';
	import { Input } from '$lib/components/ui/input';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import { getContext, hasContext } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { InputType } from './types';
	import { typeToIcon } from './utils.svelte';

	type Props = WithElementRef<Omit<HTMLInputAttributes, 'type' | 'files'> & { type?: InputType }>;
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		value = $bindable(),
		type,
		class: className,
		id,
		required,
		oninput,
		transformer = (value) => value,
		...restProps
	}: Props & {
		transformer?: (value: any) => void;
	} = $props();

	const isInvalid = $derived(required && (value === null || value === undefined || value === ''));

	const formValidator: FormValidator = getContext('FormValidator');
	$effect(() => {
		formValidator.set(id, isInvalid);
	});
</script>

<div class="relative">
	{#if type}
		<span class="absolute left-3 top-1/2 -translate-y-1/2 items-center">
			<Icon icon={hasContext('icon') ? getContext('icon') : typeToIcon[type]} />
		</span>
	{/if}

	<Input
		bind:ref
		data-slot="input-general"
		class={cn(
			'pl-9 ring-1',
			isInvalid ? 'placeholder:text-destructive/60 placeholder:text-xs' : '',
			isInvalid ? 'ring-destructive' : '',
			className
		)}
		{type}
		bind:value
		placeholder={isInvalid ? 'Required' : ''}
		oninput={(e) => {
			value = transformer(value);
			oninput?.(e);
		}}
		{...restProps}
	/>
</div>
