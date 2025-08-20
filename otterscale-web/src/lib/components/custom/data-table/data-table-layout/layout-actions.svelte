<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: DropdownMenuPrimitive.ContentProps & {
		portalProps?: DropdownMenuPrimitive.PortalProps;
	} = $props();
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<div class="w-full">
				<Button variant="ghost" size="icon" class="float-right" {...props}>
					<Icon icon="ph:dots-three" />
				</Button>
			</div>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content bind:ref class={cn('*:text-xs', className)} {...restProps}>
		{@render children?.()}
	</DropdownMenu.Content>
</DropdownMenu.Root>
