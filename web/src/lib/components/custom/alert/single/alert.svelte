<script lang="ts" module>
	import type { WithElementRef } from 'bits-ui';
	import { setContext } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';
	import { tv, type VariantProps } from 'tailwind-variants';

	import { alertVariants as baseAlertVariants } from '$lib/components/ui/alert/alert.svelte';
	import * as Alert from '$lib/components/ui/alert/index';
	import { cn } from '$lib/utils.js';

	export const alertVariants = tv({
		base: cn(baseAlertVariants.base, '[&>svg]:size-6'),
		variants: {
			variant: {
				default: '*:data-[slot=alert-description]:text-muted-foreground bg-card/10 [&>svg]:text-current',
				information:
					'text-chart-3 *:data-[slot=alert-description]:text-muted-foreground bg-chart-3/10 [&>svg]:text-current',
				warning:
					'text-chart-4 *:data-[slot=alert-description]:text-muted-foreground bg-chart-4/10 [&>svg]:text-current',
				destructive:
					'text-destructive *:data-[slot=alert-description]:text-muted-foreground bg-destructive/10 [&>svg]:text-current',
			},
		},
		defaultVariants: {
			variant: 'default',
		},
	});

	export type AlertVariant = VariantProps<typeof alertVariants>['variant'];
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		variant = 'default',
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		variant?: AlertVariant;
	} = $props();

	setContext('variant', variant);
</script>

<Alert.Root bind:ref data-slot="alert-root" class={cn(alertVariants({ variant }), className)} {...restProps}>
	{@render children?.()}
</Alert.Root>
