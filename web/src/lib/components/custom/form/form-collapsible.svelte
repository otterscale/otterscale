<script lang="ts" module>
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { Collapsible as CollapsiblePrimitive } from 'bits-ui';
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';

	let {
		ref = $bindable(null),
		class: className,
		children,
		descriptor,
		...restProps
	}: CollapsiblePrimitive.TriggerProps & {
		descriptor?: Snippet<[{ open: boolean }]>;
	} = $props();

	let open = $state(false);
</script>

<Collapsible.Root bind:ref bind:open>
	<Collapsible.Trigger
		class={cn('relative w-full px-4 text-lg font-bold', className)}
		{...restProps}
	>
		<span
			class={cn(open ? 'invisible' : 'text-muted-foreground', 'flex items-center justify-center')}
		>
			{@render descriptor?.({ open })}
		</span>
		<Button variant="ghost" size="icon" class=" absolute right-0 top-1/2 -translate-y-1/2">
			<Icon icon="ph:caret-up-down" />
		</Button>
	</Collapsible.Trigger>
	<Collapsible.Content>
		{@render children?.()}
	</Collapsible.Content>
</Collapsible.Root>
