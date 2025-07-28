<script lang="ts" module>
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
</script>

<script lang="ts" generics="TData">
	let { text = 'No Results', class: className }: { text?: string; class?: string } = $props();
</script>

<div class={cn('flex items-center justify-center gap-1 p-4', className)}>
	<Icon icon="ph:magnifying-glass" />
	{text}
</div>
