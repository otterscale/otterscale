<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { cn } from '$lib/utils';

	let {
		checked = false,
		onCheckedChange = (v) => (checked = v),
		class: className,
		...restProps
	}: ComponentProps<typeof Checkbox> = $props();
</script>

<Checkbox class={cn(className)} bind:checked={() => checked, onCheckedChange} {...restProps} />
