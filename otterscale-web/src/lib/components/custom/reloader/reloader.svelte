<script lang="ts">
	import Label from '$lib/components/ui/label/label.svelte';
	import Switch from '$lib/components/ui/switch/switch.svelte';
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

	let { reloadManager, align = 'right' }: { reloadManager: ReloadManager; align?: AlignType } =
		$props();
</script>

<div class={cn('flex items-center justify-start gap-2', getAlignClassName(align))}>
	<Label class="bg-muted flex h-9 items-center justify-center rounded-sm p-2">Refresh</Label>
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
</div>
