<script lang="ts" module>
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
	import type { AlignType } from './type';
	import type { ReloadManager } from './utils.svelte';

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
	let { reloadManager, align = 'right' }: { reloadManager: ReloadManager; align?: AlignType } =
		$props();
</script>

<div class={cn('flex items-center justify-start gap-1', getAlignClassName(align))}>
	<Switch
		bind:checked={reloadManager.state}
		onCheckedChange={() => {
			if (reloadManager.state) {
				reloadManager.restart();
			} else {
				reloadManager.stop();
			}
		}}
	/>
	{m.auto_update()}
</div>
