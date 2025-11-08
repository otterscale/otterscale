<script lang="ts" module>
	import { Collapsible as CollapsiblePrimitive } from 'bits-ui';

	import * as Collapsible from '$lib/components/ui/collapsible';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
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
		'relative flex w-full items-center gap-2 bg-primary-foreground p-2 text-card-foreground data-[state=closed]:rounded data-[state=open]:rounded-t',
		className
	)}
	{...restProps}
>
	{@render children?.()}
	<Button variant="ghost" size="icon" class="absolute top-1/2 right-0 -translate-y-1/2">
		<Icon icon="ph:caret-up-down" />
	</Button>
</Collapsible.Trigger>
