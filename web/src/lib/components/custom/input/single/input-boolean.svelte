<script lang="ts" module>
	const type = 'boolean';
</script>

<script lang="ts">
	import * as Switch from '$lib/components/ui/switch';
	import { BORDER_INPUT_CLASSNAME, typeToIcon } from './utils.svelte';
	import Icon from '@iconify/svelte';

	import { Switch as SwitchPrimitive } from 'bits-ui';
	import { cn, type WithoutChildrenOrChild } from '$lib/utils.js';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	let {
		ref = $bindable(null),
		class: className,
		value: checked = $bindable(false),
		...restProps
	}: WithoutChildrenOrChild<SwitchPrimitive.RootProps> = $props();
</script>

<div class={cn(BORDER_INPUT_CLASSNAME, 'justify-between')}>
	<span class="flex items-center gap-2">
		<span class="pl-3">
			<Icon icon={typeToIcon[type]} />
		</span>
		<Badge variant={checked ? 'default' : 'destructive'}
			>{String(checked).toLocaleUpperCase()}</Badge
		>
	</span>
	<Switch.Root
		bind:ref
		bind:checked
		data-slot="input-boolean"
		class={cn('my-2 mr-3 cursor-pointer', className)}
		{...restProps}
	/>
</div>
