<script lang="ts">
	import {
		BORDER_INPUT_CLASSNAME,
		UNFOCUS_INPUT_CLASSNAME,
		typeToIcon,
		PasswordManager
	} from './utils.svelte';
	import InputRequired from './input-required.svelte';

	import Icon from '@iconify/svelte';
	import { Input } from '$lib/components/ui/input';

	import type { HTMLInputAttributes } from 'svelte/elements';
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';

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

{#if isNotFilled}
	<InputRequired {isNotFilled} />
{/if}
<div class={cn(BORDER_INPUT_CLASSNAME)}>
	<span class="pl-3">
		<Icon icon={typeToIcon['password']} />
	</span>
	<Input
		bind:ref
		data-slot="input-password"
		class={cn(UNFOCUS_INPUT_CLASSNAME, className)}
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
