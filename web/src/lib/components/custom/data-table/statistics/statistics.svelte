<script lang="ts" module>
	export type Type = 'data' | 'count' | 'ratio';
	export type TypeAccessor = { value: Type };
</script>

<script lang="ts">
	import { setContext } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Card from '$lib/components/ui/card';
	import { cn, type WithElementRef } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		children,
		type,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		type: Type;
	} = $props();

	const typeAccessor: TypeAccessor = {
		get value() {
			return type;
		},
		set value(newType: Type) {
			type = newType;
		}
	};
	setContext('typeAccessor', typeAccessor);
</script>

<Card.Root
	bind:ref
	data-slot="data-table-statistics"
	class={cn('relative overflow-hidden', className)}
	{...restProps}
>
	{@render children?.()}
</Card.Root>
