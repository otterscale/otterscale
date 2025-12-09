<script lang="ts" module>
	import { Switch as SwitchPrimitive } from 'bits-ui';

	import { Switch } from '$lib/components/ui/switch/index.js';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
	import { type WithoutChildrenOrChild } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: WithoutChildrenOrChild<SwitchPrimitive.RootProps> = $props();
</script>

<div class={cn('flex items-center justify-start gap-1', className)}>
	<Switch bind:ref {...restProps} />
	<span class="text-sm font-medium capitalize">{m.auto_update()}</span>
</div>
