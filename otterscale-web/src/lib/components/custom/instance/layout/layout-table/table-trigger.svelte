<script lang="ts">
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { Collapsible as CollapsiblePrimitive } from 'bits-ui';
	import { cn } from '$lib/utils';

	import Icon from '@iconify/svelte';
	import Button from '$lib/components/ui/button/button.svelte';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: CollapsiblePrimitive.TriggerProps = $props();
</script>

<Collapsible.Trigger
	bind:ref
	data-slot="table-trigger"
	class={cn(
		'bg-primary-foreground text-card-foreground relative flex w-full items-center gap-2 p-2 data-[state=closed]:rounded data-[state=open]:rounded-t',
		className
	)}
	{...restProps}
>
	{@render children?.()}
	<Button variant="ghost" size="icon" class="absolute top-1/2 right-0 -translate-y-1/2">
		<Icon icon="ph:caret-up-down" />
	</Button>
</Collapsible.Trigger>
