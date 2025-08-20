<script lang="ts">
	import Input from '$lib/components/ui/input/input.svelte';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import type { FilterManager } from './utils';
	import { m } from '$lib/paraglide/messages';

	let { filterManager }: { filterManager: FilterManager } = $props();
</script>

<span class="relative h-8">
	<Input
		placeholder={m.applications_store_chart_filter_name()}
		type="text"
		bind:value={filterManager.searchedName}
		class="h-8 w-40 pr-9 pl-9"
	/>
	<Icon icon="ph:magnifying-glass" class="absolute top-1/2 left-3 -translate-y-1/2" />
	<button
		onclick={() => {
			filterManager.resetName();
		}}
	>
		<Icon
			icon="ph:x"
			class={cn(
				filterManager.searchedName ? 'visible' : 'hidden',
				'absolute top-1/2 right-3 -translate-y-1/2'
			)}
		/>
	</button>
</span>
