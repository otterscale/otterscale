<script lang="ts" module>
	import { Button, buttonVariants, type ButtonProps } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		href = undefined,
		type = 'button',
		disabled,
		variant = 'default',
		size = 'sm',
		children,
		...restProps
	}: ButtonProps = $props();
</script>

<Button
	bind:ref
	data-slot="select-actions"
	class={cn('cursor-pointer', buttonVariants({ variant, size }), className)}
	{href}
	{type}
	{disabled}
	{...restProps}
>
	{@render children?.()}
</Button>
