<script lang="ts">
	import * as Command from '$lib/components/ui/command';
	import { cn } from '$lib/utils.js';
	import { Command as CommandPrimitive } from 'bits-ui';

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: CommandPrimitive.SeparatorProps = $props();
</script>

<Command.Separator bind:ref data-slot="select-separator" class={cn(className)} {...restProps} />
