<script lang="ts">
	import * as Command from '$lib/components/ui/command';

	import { Command as CommandPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: CommandPrimitive.SeparatorProps = $props();
</script>

<Command.Separator bind:ref data-slot="select-separator" class={cn(className)} {...restProps} />
