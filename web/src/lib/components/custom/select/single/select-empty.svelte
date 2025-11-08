<script lang="ts" module>
	import { Command as CommandPrimitive } from 'bits-ui';

	import * as Command from '$lib/components/ui/command';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: CommandPrimitive.EmptyProps = $props();
</script>

<Command.Empty bind:ref data-slot="select-empty" class={cn(className)} {...restProps} />
