<script lang="ts" module>
	import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';
	import { FilterManager } from './utils';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let { filterManager }: { filterManager: FilterManager } = $props();
</script>

<Button
	variant="secondary"
	class={cn('flex h-8 items-center gap-2', filterManager.isFiltered ? 'visible' : 'hidden')}
	onclick={() => {
		filterManager.reset();
	}}
>
	<Icon icon="ph:arrow-counter-clockwise" class="h-3 w-3" />
	Reset
</Button>
