<script lang="ts" module>
	import { Menubar as MenubarPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';

	import type { Setter } from './types';

	import * as Menubar from '$lib/components/ui/menubar/index';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		// inset,
		value,
		onclick,
		...restProps
	}: MenubarPrimitive.ItemProps & {
		// inset?: boolean;
		value?: string;
	} = $props();

	const setter: Setter = getContext('setter');
</script>

<Menubar.Item
	bind:ref
	class={cn('p-1 text-xs hover:cursor-pointer [&_svg]:m-1', className)}
	{...restProps}
	onclick={(e) => {
		setter(value);
		onclick?.(e);
	}}
/>
