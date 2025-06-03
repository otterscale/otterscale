<script lang="ts" module>
	const type = 'password';
</script>

<script lang="ts">
	import {
		BORDER_INPUT_CLASSNAME,
		UNFOCUS_INPUT_CLASSNAME,
		typeToIcon,
		PasswordManager
	} from './utils.svelte';

	import Icon from '@iconify/svelte';
	import { Input } from '$lib/components/ui/input';

	import type { HTMLInputAttributes } from 'svelte/elements';
	import { cn, type WithElementRef } from '$lib/utils.js';

	type Props = WithElementRef<Omit<HTMLInputAttributes, 'type'> & { type?: 'password' }>;

	let {
		ref = $bindable(null),
		value = $bindable(),
		class: className,
		...restProps
	}: Props = $props();

	const { files, ...restProperties } = restProps;
	const passwordManager = new PasswordManager();
</script>

<div class={cn(BORDER_INPUT_CLASSNAME)}>
	<span class="pl-3">
		<Icon icon={typeToIcon[type]} />
	</span>
	<Input
		bind:ref
		data-slot="input-password"
		class={cn(UNFOCUS_INPUT_CLASSNAME, className)}
		type={passwordManager.isVisible ? 'text' : 'password'}
		bind:value
		{...restProperties}
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
