<script lang="ts">
	import { getContext } from 'svelte';
	import Icon from '@iconify/svelte';
	import { Badge, type BadgeVariant } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';

	import { typeToIcon, ValuesManager, InputManager } from './utils.svelte';

	import type { HTMLAnchorAttributes } from 'svelte/elements';
	import { cn, type WithElementRef } from '$lib/utils';

	let {
		ref = $bindable(null),
		href,
		class: className,
		variant = 'outline',
		disabled,
		children,
		...restProps
	}: WithElementRef<HTMLAnchorAttributes> & {
		variant?: BadgeVariant;
		disabled?: boolean;
	} = $props();

	const inputManager: InputManager = getContext('InputManager');
	const valuesManager: ValuesManager = getContext('ValuesManager');
</script>

{#if valuesManager.values.length > 0}
	<div class="flex flex-wrap gap-1">
		{#each valuesManager.values as value}
			<Badge
				{href}
				bind:ref
				data-slot="input-viewer"
				class={cn('flex h-6 items-center gap-2', className)}
				{...restProps}
				{variant}
			>
				<span class="flex items-center gap-1">
					<Icon icon={typeToIcon[inputManager.type]} />
					{value}
				</span>
				<Button
					class="size-3 cursor-pointer"
					aria-label="Remove"
					size="icon"
					variant="ghost"
					{disabled}
					onclick={(e) => {
						valuesManager.remove(value);
					}}
				>
					<Icon icon="ph:x-circle" class="size-3" />
				</Button>
			</Badge>
		{/each}
	</div>
{/if}
