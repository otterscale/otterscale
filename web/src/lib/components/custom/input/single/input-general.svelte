<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import { getContext, hasContext } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';

	import { Input } from '$lib/components/ui/input';
	import { cn } from '$lib/utils.js';

	import type { InputType } from './types';
	import { typeToIcon } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		value = $bindable(),
		type = 'text',
		class: className,
		required,
		oninput,
		transformer = (value) => value,
		validator = () => true,
		invalid = $bindable(),
		...restProps
	}: WithElementRef<
		Omit<HTMLInputAttributes, 'type' | 'files'> & { type?: InputType | undefined }
	> & {
		transformer?: (value: any) => void;
		validator?: (value?: any) => boolean;
		invalid?: boolean | null | undefined;
	} = $props();

	const isNull = $derived(required && (value === null || value === undefined || value === ''));
	const isValidated = $derived(value ? validator(value) : true);
	const isInvalid = $derived(isNull || !isValidated);
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
				? 'text-destructive/60 placeholder:text-xs placeholder:text-destructive/60 focus:placeholder:invisible'
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
