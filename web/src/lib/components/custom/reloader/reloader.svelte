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
		checked = $bindable(false),
		...restProps
	}: WithoutChildrenOrChild<SwitchPrimitive.RootProps> = $props();
</script>

<div class={cn('flex items-center justify-start gap-1', className)}>
	<Switch bind:ref bind:checked {...restProps} />
	{m.auto_update()}
</div>
