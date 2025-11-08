<script lang="ts" module>
	import type { WithElementRef } from 'bits-ui';
	import { setContext } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';
	import { tv, type VariantProps } from 'tailwind-variants';

	import type { AlertType } from './types';

	import { alertVariants as baseAlertVariants } from '$lib/components/ui/alert/alert.svelte';
	import * as Alert from '$lib/components/ui/alert/index';
	import { cn } from '$lib/utils.js';

	export const alertVariants = tv({
		base: cn(baseAlertVariants.base, '[&>svg]:size-6'),
		variants: {
			variant: {
				default:
					'*:data-[slot=alert-description]:text-muted-foreground bg-card/10 [&>svg]:text-current',
				destructive:
					'border-destructive/50 text-destructive *:data-[slot=alert-description]:text-destructive/90 [&>svg]:text-current',
				success:
					'bg-card border-green-800/50 text-green-700 *:data-[slot=alert-description]:text-green-600 dark:border-green-800 dark:text-green-500 dark:*:data-[slot=alert-description]:text-green-600 [&>svg]:text-green-800',
				warning:
					'bg-card border-yellow-500/50 text-yellow-700 *:data-[slot=alert-description]:text-yellow-600 dark:border-yellow-500 dark:text-yellow-400 dark:*:data-[slot=alert-description]:text-yellow-700 [&>svg]:text-current',
				information:
					'bg-card border-blue-500/50 text-blue-700 *:data-[slot=alert-description]:text-blue-500 dark:border-blue-500 dark:text-blue-400 dark:*:data-[slot=alert-description]:text-blue-600 [&>svg]:text-blue-500'
			}
		},
		defaultVariants: {
			variant: 'default'
		}
	});

	export type AlertVariant = VariantProps<typeof alertVariants>['variant'];
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		alert,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		alert: AlertType;
	} = $props();

	setContext('variant', alert.variant);
</script>

<Alert.Root
	bind:ref
	data-slot="alert-root"
	class={cn(alertVariants({ variant: alert.variant }), className)}
	{...restProps}
>
	{@render children?.()}
</Alert.Root>
