<script lang="ts" module>
	import { tv, type VariantProps } from 'tailwind-variants';

	export const tableVariants = tv({
		variants: {
			variant: {
				outline: 'border-border border shadow-xs rounded-lg',
				ghost: 'hover:bg-accent/30 rounded-lg'
			}
		},
		defaultVariants: {
			variant: 'outline'
		}
	});

	export type TableVariant = VariantProps<typeof tableVariants>['variant'];
</script>

<script lang="ts">
	import type { HTMLTableAttributes } from 'svelte/elements';

	import * as Table from '$lib/components/ui/table';
	import { cn, type WithElementRef } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		variant,
		children,
		...restProps
	}: WithElementRef<HTMLTableAttributes> & { variant?: TableVariant } = $props();
</script>

<div class={tableVariants({ variant })}>
	<Table.Root bind:ref data-slot="table" class={cn(className)} {...restProps}>
		{@render children?.()}
	</Table.Root>
</div>
