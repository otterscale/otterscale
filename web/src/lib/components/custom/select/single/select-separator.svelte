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
	}: CommandPrimitive.SeparatorProps = $props();
</script>

<Command.Separator bind:ref data-slot="select-separator" class={cn(className)} {...restProps} />
