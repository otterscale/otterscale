<script lang="ts" module>
	import { Button, type Props } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		class: className,
		variant = 'outline',
		size = 'sm',
		ref = $bindable(null),
		href = undefined,
		type = 'button',
		disabled,
		children,
		...restProps
	}: Props = $props();
</script>

<div class="flex h-full items-end justify-end">
	<Button
		bind:ref
		data-slot="alert-action"
		class={cn('hover:cursor-pointer', className)}
		{variant}
		{size}
		{type}
		{disabled}
		{...restProps}
	>
		{@render children?.()}
	</Button>
</div>
