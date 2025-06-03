<script lang="ts" module>
	const type = 'color';
</script>

<script lang="ts">
	import { BORDER_INPUT_CLASSNAME, UNFOCUS_INPUT_CLASSNAME, typeToIcon } from './utils.svelte';

	import Icon from '@iconify/svelte';
	import { Input } from '$lib/components/ui/input';

	import type { HTMLInputAttributes, HTMLInputTypeAttribute } from 'svelte/elements';
	import { cn, type WithElementRef } from '$lib/utils.js';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	type Props = WithElementRef<
		Omit<HTMLInputAttributes, 'type'> & {
			type?: Exclude<HTMLInputTypeAttribute, 'file' | 'password'>;
		}
	>;

	let {
		ref = $bindable(null),
		value = $bindable(),
		class: className,
		...restProps
	}: Props = $props();

	const { files, ...restProperties } = restProps;
</script>

<div class={cn(BORDER_INPUT_CLASSNAME, 'h-10 justify-between')}>
	<span class="flex items-center gap-2">
		<span class="flex items-center gap-2">
			<span class="pl-3">
				<Icon icon={typeToIcon[type]} />
			</span>
			<Badge variant="outline">{value}</Badge>
		</span>
	</span>
	<Input
		bind:ref
		data-slot="input-color"
		class={cn(
			UNFOCUS_INPUT_CLASSNAME,
			'mr-3 aspect-square h-7 w-fit cursor-pointer p-0',
			className
		)}
		{type}
		bind:value
		{...restProperties}
	/>
</div>
