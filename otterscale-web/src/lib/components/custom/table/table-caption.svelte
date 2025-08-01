<script lang="ts" module>
	import * as Table from '$lib/components/ui/table';
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';

	type Align = 'top' | 'bottom';

	function getCaptionAlignClassName(align: Align) {
		if (align === 'top') {
			return 'caption-top';
		} else if (align === 'bottom') {
			return 'caption-bottom';
		}
	}
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		align = 'top',
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLElement>> & {
		align?: Align;
	} = $props();
</script>

<Table.Caption
	bind:ref
	data-slot="table-caption"
	class={cn(getCaptionAlignClassName(align), className)}
	{...restProps}
>
	{@render children?.()}
</Table.Caption>
