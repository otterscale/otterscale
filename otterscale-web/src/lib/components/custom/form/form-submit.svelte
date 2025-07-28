<script lang="ts" module>
	import { FormValidator } from '$lib/components/custom/form';
	import { Button, buttonVariants, type ButtonProps } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils.js';
	import { getContext } from 'svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		href = undefined,
		disabled,
		variant = 'default',
		size = 'sm',
		children,
		...restProps
	}: ButtonProps & {} = $props();

	const formValidator: FormValidator = getContext('FormValidator');
</script>

<Button
	bind:ref
	data-slot="form-submit"
	class={cn('w-fit cursor-pointer shadow', buttonVariants({ variant, size }), className)}
	{href}
	type="submit"
	disabled={formValidator.isInvalid}
	{...restProps}
>
	{@render children?.()}
</Button>
