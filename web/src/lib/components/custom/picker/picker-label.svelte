<script lang="ts" module>
	import Label from '$lib/components/ui/label/label.svelte';
	import { cn } from '$lib/utils.js';
	import { Label as LabelPrimitive } from 'bits-ui';
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
		'bg-muted flex items-center justify-center whitespace-nowrap rounded-lg p-2',
		className
	)}
>
	{@render children?.()}
</Label>
