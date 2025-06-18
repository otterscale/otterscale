<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import * as Tooptip from '$lib/components/ui/tooltip';
	import { Tooltip, type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';
	import { enabled } from './data';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {} = $props();
</script>

{#if enabled}
	{@render children?.()}
{:else}
	<Tooltip.Provider>
		<Tooltip.Root>
			<Tooltip.Trigger>
				{@render children?.()}
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
