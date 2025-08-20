<script lang="ts">
	import * as Command from '$lib/components/ui/command';
	import { cn } from '$lib/utils.js';
	import { Command as CommandPrimitive } from 'bits-ui';

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: CommandPrimitive.EmptyProps = $props();
</script>

<Command.Empty bind:ref data-slot="select-empty" class={cn(className)} {...restProps} />
