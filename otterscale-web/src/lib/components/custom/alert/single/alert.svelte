<script lang="ts">
	import { setContext } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';
	import { tv, type VariantProps } from 'tailwind-variants';

	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';

	import { alertVariants as baseAlertVariants } from '$lib/components/ui/alert/alert.svelte';
	import * as Alert from '$lib/components/ui/alert/index';

	export const alertVariants = tv({
		base: cn(baseAlertVariants.base, '[&>svg]:size-6'),
		variants: {
			variant: {
				default:
					'*:data-[slot=alert-description]:text-muted-foreground [&>svg]:text-current bg-card/10',
				information:
					'text-chart-3 *:data-[slot=alert-description]:text-muted-foreground  [&>svg]:text-current bg-chart-3/10',
				warning:
					'text-chart-4 *:data-[slot=alert-description]:text-muted-foreground [&>svg]:text-current bg-chart-4/10',
				destructive:
					'text-destructive *:data-[slot=alert-description]:text-muted-foreground [&>svg]:text-current bg-destructive/10'
			}
		},
		defaultVariants: {
			variant: 'default'
		}
	});

	export type AlertVariant = VariantProps<typeof alertVariants>['variant'];

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

<Alert.Root
	bind:ref
	data-slot="alert-root"
	class={cn(alertVariants({ variant }), className)}
	{...restProps}
>
	{@render children?.()}
</Alert.Root>
