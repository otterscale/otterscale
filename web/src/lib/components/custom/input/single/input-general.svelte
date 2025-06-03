<script lang="ts">
	import { BORDER_INPUT_CLASSNAME, UNFOCUS_INPUT_CLASSNAME, typeToIcon } from './utils.svelte';

	import Icon from '@iconify/svelte';
	import { Input } from '$lib/components/ui/input';

	import type { HTMLInputAttributes, HTMLInputTypeAttribute } from 'svelte/elements';
	import { cn, type WithElementRef } from '$lib/utils.js';

	type Props = WithElementRef<
		Omit<HTMLInputAttributes, 'type'> & {
			type?: Exclude<HTMLInputTypeAttribute, 'file' | 'password'>;
		}
	>;

	let {
		ref = $bindable(null),
		value = $bindable(),
		type,
		class: className,
		...restProps
	}: Props = $props();

	const { files, ...restProperties } = restProps;
</script>

<div class={cn(BORDER_INPUT_CLASSNAME)}>
	{#if type}
		<span class="pl-3">
			<Icon icon={typeToIcon[type]} />
		</span>
	{/if}
	<Input
		bind:ref
		data-slot="input-general"
		class={cn(UNFOCUS_INPUT_CLASSNAME, className)}
		{type}
		bind:value
		{...restProperties}
	/>
</div>
