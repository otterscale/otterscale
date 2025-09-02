<script lang="ts" module>
	import * as Command from '$lib/components/ui/command';
	import { cn } from '$lib/utils.js';
	import { Command as CommandPrimitive } from 'bits-ui';
</script>

<script lang="ts">
	let { ref = $bindable(null), class: className, ...restProps }: CommandPrimitive.ListProps = $props();
</script>

<Command.List bind:ref data-slot="select-list" class={cn(className)} {...restProps} />
