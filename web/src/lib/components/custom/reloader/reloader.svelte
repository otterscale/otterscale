<script lang="ts" module>
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
	import { type WithoutChildrenOrChild } from '$lib/utils.js';
	import { Switch as SwitchPrimitive } from 'bits-ui';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		checked = $bindable(false),
		...restProps
	}: WithoutChildrenOrChild<SwitchPrimitive.RootProps> = $props();
</script>

<div class={cn('flex items-center justify-start gap-1')}>
	<Switch bind:ref bind:checked {...restProps} />
	{m.auto_update()}
</div>
