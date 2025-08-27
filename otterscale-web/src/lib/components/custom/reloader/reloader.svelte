<script lang="ts" module>
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import type { AlignType } from './type';
	import type { ReloadManager } from './utils.svelte';
	import { fade } from 'svelte/transition';

	function getAlignClassName(align: AlignType) {
		if (align === 'left') {
			return 'flex justify-start';
		} else if (align === 'right') {
			return 'flex justify-end';
		} else {
			return '';
		}
	}
</script>

<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';

	let { reloadManager, align = 'right' }: { reloadManager: ReloadManager; align?: AlignType } =
		$props();
</script>

<div class={cn('flex items-center justify-start gap-1', getAlignClassName(align))}>
	<Button
		onclick={() => {
			reloadManager.state = !reloadManager.state;
			if (reloadManager.state) {
				reloadManager.restart();
			} else {
				reloadManager.stop();
			}
		}}
	>
		{#if reloadManager.state}
			<Icon
				icon="ph:arrows-clockwise"
				class="animate-spin opacity-100 transition-opacity duration-300"
			/>
		{/if}
		{m.auto_refresh()}
	</Button>
</div>
