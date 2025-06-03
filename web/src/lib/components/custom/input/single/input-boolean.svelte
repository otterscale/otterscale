<script lang="ts" module>
	const type = 'boolean';
</script>

<script lang="ts">
	import * as Switch from '$lib/components/ui/switch';
	import { BORDER_INPUT_CLASSNAME, typeToIcon } from './utils.svelte';
	import Icon from '@iconify/svelte';
	import { z, type ZodFirstPartySchemaTypes } from 'zod';
	import { InputValidator } from './utils.svelte';

	import { Switch as SwitchPrimitive, type WithoutChildrenOrChild } from 'bits-ui';
	import { cn } from '$lib/utils.js';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	let {
		ref = $bindable(null),
		class: className,
		required,
		value: checked = $bindable(false),
		...restProps
	}: WithoutChildrenOrChild<SwitchPrimitive.RootProps> & {} = $props();
</script>

<div
	class={cn(
		BORDER_INPUT_CLASSNAME,
		checked !== true && checked !== false ? 'ring-destructive ring-1' : '',
		'justify-between',
		className
	)}
>
	<span class="flex items-center gap-2">
		<span class="pl-3">
			<Icon icon={typeToIcon[type]} />
		</span>
		{#if checked === true}
			<Badge variant="default">True</Badge>
		{:else if checked === false}
			<Badge variant="outline">False</Badge>
		{:else if checked === null || checked === undefined}
			<Badge variant="secondary">null</Badge>
		{:else}
			<Badge variant="destructive">Invalid</Badge>
		{/if}
	</span>
	<Switch.Root
		bind:ref
		bind:checked
		data-slot="input-boolean"
		class={cn('my-2 mr-3 cursor-pointer')}
		{...restProps}
	/>
</div>
<div class="transition-all duration-500">
	{#if checked === null || checked === undefined}
		<div class="animate-in fade-in flex items-center gap-1">
			<Icon icon="ph:asterisk" class="text-destructive size-2" />
			<p class="text-destructive text-xs">Required</p>
		</div>
	{/if}
</div>
