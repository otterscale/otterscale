<script lang="ts" module>
	import { Input } from '$lib/components/ui/input';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLInputAttributes } from 'svelte/elements';
</script>

<script lang="ts">
	import {
		BORDER_INPUT_CLASSNAME,
		PasswordManager,
		RING_INVALID_INPUT_CLASSNAME,
		RING_VALID_INPUT_CLASSNAME,
		typeToIcon,
		UNFOCUS_INPUT_CLASSNAME
	} from './utils.svelte';

	let {
		ref = $bindable(null),
		value = $bindable(),
		required,
		class: className,
		...restProps
	}: WithElementRef<Exclude<HTMLInputAttributes, 'type'>> & { type?: 'password' } = $props();

	const passwordManager = new PasswordManager();
	const isNotFilled = $derived(required && !value);
</script>

<div
	class={cn(
		BORDER_INPUT_CLASSNAME,
		isNotFilled ? RING_INVALID_INPUT_CLASSNAME : RING_VALID_INPUT_CLASSNAME
	)}
>
	<span class="pl-3">
		<Icon icon={typeToIcon['password']} />
	</span>

	<Input
		bind:ref
		data-slot="input-password"
		placeholder={isNotFilled ? 'Required' : ''}
		class={cn(
			UNFOCUS_INPUT_CLASSNAME,
			isNotFilled ? 'placeholder:text-destructive/60 placeholder:text-xs' : '',
			className
		)}
		type={passwordManager.isVisible ? 'text' : 'password'}
		bind:value
		{...restProps}
	/>

	<button
		type="button"
		class="pr-3 hover:cursor-pointer focus:outline-none"
		onmousedown={() => {
			passwordManager.enable();
		}}
		onmouseup={() => {
			passwordManager.disable();
		}}
	>
		<Icon icon={passwordManager.isVisible ? 'ph:eye' : 'ph:eye-slash'} />
	</button>
</div>
