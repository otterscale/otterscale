<script lang="ts" module>
	import { Badge, type BadgeVariant } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import { getContext, hasContext } from 'svelte';
	import type { HTMLAnchorAttributes } from 'svelte/elements';
</script>

<script lang="ts">
	import { InputManager, typeToIcon, ValuesManager } from './utils.svelte';

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
			{#if !(value == undefined || value == null)}
				<Badge
					{href}
					bind:ref
					data-slot="input-viewer"
					class={cn('flex h-6 items-center gap-3', className)}
					{...restProps}
					{variant}
				>
					<span class="flex items-center gap-1">
						<Icon icon={hasContext('icon') ? getContext('icon') : typeToIcon[inputManager.type]} />
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
						<Icon icon="ph:x-circle" class="text-muted-foreground" />
					</Button>
				</Badge>
			{/if}
		{/each}
	</div>
{/if}
