<script lang="ts">
	import Icon from '@iconify/svelte';
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';

	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		open = $bindable(),
		class: className,
		children,
		...restProps
	}: DropdownMenuPrimitive.ContentProps & {
		open?: boolean;
		portalProps?: DropdownMenuPrimitive.PortalProps;
	} = $props();
</script>

<DropdownMenu.Root bind:open>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<div class="w-full">
				<Button variant="ghost" size="icon" class="float-right" {...props}>
					<Icon icon="ph:dots-three" />
				</Button>
			</div>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content bind:ref class={cn(className)} {...restProps}>
		{@render children?.()}
	</DropdownMenu.Content>
</DropdownMenu.Root>
