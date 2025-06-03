<script lang="ts">
	import * as Command from '$lib/components/ui/command';

	import { Command as CommandPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: CommandPrimitive.EmptyProps = $props();
</script>

<Command.Empty bind:ref data-slot="select-empty" class={cn(className)} {...restProps} />
