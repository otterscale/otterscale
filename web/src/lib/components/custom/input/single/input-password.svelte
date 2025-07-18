<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import { General } from '.';
</script>

<script lang="ts">
	import { PasswordManager } from './utils.svelte';

	let {
		id,
		ref = $bindable(null),
		value = $bindable(),
		required,
		class: className,
		...restProps
	}: WithElementRef<Exclude<HTMLInputAttributes, 'type'>> & { type?: 'password' } = $props();

	const passwordManager = new PasswordManager();
	const isNotFilled = $derived(required && !value);
</script>

<div class="relative">
	<General
		bind:ref
		data-slot="input-password"
		class="pr-9"
		type={passwordManager.isVisible ? 'text' : 'password'}
		{required}
		aria-invalid={isNotFilled}
		bind:value
		{...restProps}
	/>

	<button
		type="button"
		class="absolute right-3 top-1/2 -translate-y-1/2 items-center hover:cursor-pointer focus:outline-none"
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
