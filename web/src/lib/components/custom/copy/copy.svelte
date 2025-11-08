<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { scale } from 'svelte/transition';

	import { Button } from '$lib/components/ui/button';
	import { UseClipboard } from '$lib/hooks/use-clipboard.svelte';
	import { cn } from '$lib/utils';

	import type { CopyProps } from './type.d';
</script>

<script lang="ts">
	let {
		class: className,
		variant = 'ghost',
		size = 'icon',
		ref = $bindable(null),
		text,
		onCopy,
		...restProps
	}: CopyProps = $props();

	const clipboard = new UseClipboard();
</script>

<Button
	bind:ref
	data-slot="copy"
	{variant}
	{size}
	class={cn('flex items-center gap-2', className)}
	type="button"
	onclick={async () => {
		const status = await clipboard.copy(text);
		onCopy?.(status);
	}}
	{...restProps}
>
	{@const duration = 500}
	{@const start = 0.85}
	{#if clipboard.status === 'success'}
		<span in:scale={{ duration, start }}>
			<Icon icon="ph:check" />
			<span class="sr-only">Copied</span>
		</span>
	{:else if clipboard.status === 'failure'}
		<span in:scale={{ duration, start }}>
			<Icon icon="ph:x" />
			<span class="sr-only">Failed to copy</span>
		</span>
	{:else}
		<span in:scale={{ duration, start }}>
			<Icon icon="ph:copy" />
			<span class="sr-only">Copy</span>
		</span>
	{/if}
</Button>
