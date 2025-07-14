<script lang="ts" module>
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { Tooltip, type WithElementRef, type WithoutChildren } from 'bits-ui';
	import { getContext, type Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		children,
		class: className,
		...restProps
	}: WithoutChildren<WithElementRef<HTMLAttributes<HTMLDivElement>>> & {
		children?: Snippet<[{ disable: boolean }]>;
	} = $props();

	const getIsEnterprise: () => boolean = getContext('getIsEnterprise');
</script>

{#if getIsEnterprise() === true}
	{@render children?.({ disable: !getIsEnterprise() })}
{:else}
	<Tooltip.Provider>
		<Tooltip.Root>
			<Tooltip.Trigger class={cn(className)}>
				{@render children?.({ disable: !getIsEnterprise() })}
			</Tooltip.Trigger>
			<Tooltip.Content class="bg-popover flex items-center gap-1 border p-2 shadow">
				<Icon icon="ph:info" />
				<p class="text-popover-foreground text-xs">
					This feature is available in the Enterprise edition. Contact us to learn more about
					upgrading your solution.
				</p>
			</Tooltip.Content>
		</Tooltip.Root>
	</Tooltip.Provider>
{/if}
