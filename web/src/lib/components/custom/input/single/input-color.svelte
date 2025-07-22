<script lang="ts" module>
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { Input } from '$lib/components/ui/input';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import { z, type ZodFirstPartySchemaTypes } from 'zod';
	import { BORDER_INPUT_CLASSNAME, typeToIcon } from './utils.svelte';
</script>

<script lang="ts">
	let {
		id,
		ref = $bindable(null),
		value = $bindable(),
		schema = z.string().regex(/^#[0-9a-fA-F]{6}$/),
		class: className,
		...restProps
	}: WithElementRef<Exclude<HTMLInputAttributes, 'type'>> & { type?: 'color' } & {
		schema?: ZodFirstPartySchemaTypes;
	} = $props();
</script>

<div class={cn(BORDER_INPUT_CLASSNAME, 'justify-between ring-1', className)}>
	<span class="flex items-center gap-2">
		<span class="pl-3">
			<Icon icon={typeToIcon['color']} />
		</span>
		<Badge variant="outline">{value}</Badge>
	</span>

	<Input
		bind:ref
		data-slot="input-color"
		class="mr-3 aspect-square h-7 w-fit cursor-pointer border-none p-0 shadow-none"
		type="color"
		bind:value
		{...restProps}
	/>
</div>
