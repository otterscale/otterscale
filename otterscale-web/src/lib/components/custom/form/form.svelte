<script lang="ts" module>
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';
	import { setContext } from 'svelte';
	import type { HTMLFormAttributes } from 'svelte/elements';
	import { FormValidator } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		children,
		invalid = $bindable(),
		...restProps
	}: WithElementRef<HTMLFormAttributes> & { invalid?: boolean | undefined } = $props();

	const formValidator = new FormValidator();
	setContext('FormValidator', formValidator);
	$effect(() => {
		invalid = formValidator.isInvalid;
	});
</script>

<form
	bind:this={ref}
	data-slot="form-root"
	class={cn('flex max-h-[77vh] flex-col gap-4 overflow-auto rounded-md', className)}
	{...restProps}
>
	{@render children?.()}
</form>
