<script lang="ts" module>
	import * as Command from '$lib/components/ui/command';
	import { cn } from '$lib/utils.js';
	import { Command as CommandPrimitive } from 'bits-ui';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		value = $bindable(''),
		class: className,
		...restProps
	}: CommandPrimitive.RootProps = $props();
</script>

<Command.Root bind:value bind:ref data-slot="select-option" class={cn(className)} {...restProps} />
