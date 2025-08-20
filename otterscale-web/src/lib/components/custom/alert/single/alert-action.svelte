<script lang="ts">
	import { Button, type Props } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';

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

<div class="col-start-2 flex h-full items-end justify-end">
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
