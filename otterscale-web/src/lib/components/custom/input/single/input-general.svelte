<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import { getContext, hasContext } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { InputType } from './types';
	import { typeToIcon } from './utils.svelte';

	let {
		ref = $bindable(null),
		value = $bindable(),
		type = 'text',
		class: className,
		id,
		required,
		oninput,
		transformer = (value) => value,
		invalid = $bindable(),
		...restProps
	}: WithElementRef<
		Omit<HTMLInputAttributes, 'type' | 'files'> & { type?: InputType | undefined }
	> & {
		transformer?: (value: any) => void;
		invalid?: boolean | null | undefined;
	} = $props();

	const isInvalid = $derived(required && (value === null || value === undefined || value === ''));
	$effect(() => {
		invalid = isInvalid;
	});
</script>

<div class="relative">
	{#if type}
		<span class="absolute top-1/2 left-3 -translate-y-1/2 items-center">
			<Icon icon={hasContext('icon') ? getContext('icon') : typeToIcon[type]} />
		</span>
	{/if}

	<Input
		bind:ref
		data-slot="input-general"
		class={cn(
			'pl-9 ring-1',
			isInvalid
				? 'placeholder:text-destructive/60 placeholder:text-xs focus:placeholder:invisible'
				: '',
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
