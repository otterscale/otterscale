<script lang="ts">
	import { Avatar as AvatarPrimitive } from 'bits-ui';

	import { Root } from '$lib/components/ui/avatar';
	import { cn } from '$lib/utils';

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: AvatarPrimitive.RootProps = $props();
</script>

<Root
	bind:ref
	class={cn("order-1 group-data-[variant='sent']/chat-bubble:order-2", className)}
	{...restProps}
/>
