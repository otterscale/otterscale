<script lang="ts" module>
	import Label from '$lib/components/ui/label/label.svelte';
	import Switch from '$lib/components/ui/switch/switch.svelte';
	import type { ReloadManager } from './utils.svelte';
</script>

<script lang="ts">
	let { reloadManager }: { reloadManager: ReloadManager } = $props();
</script>

<div class="flex items-center justify-start gap-2">
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
