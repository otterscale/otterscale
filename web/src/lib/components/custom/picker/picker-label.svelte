<script lang="ts" module>
	import { Label as LabelPrimitive } from 'bits-ui';

	import Label from '$lib/components/ui/label/label.svelte';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: LabelPrimitive.RootProps = $props();
</script>

<Label
	bind:ref
	{...restProps}
	class={cn(
		'flex items-center justify-center rounded-lg bg-muted p-2 whitespace-nowrap',
		className
	)}
>
	{@render children?.()}
</Label>
