<script lang="ts">
	import * as Command from '$lib/components/ui/command';

	import { Command as CommandPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: CommandPrimitive.ListProps = $props();
</script>

<Command.List bind:ref data-slot="select-list" class={cn(className)} {...restProps} />
