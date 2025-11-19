<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { VirtualMachine } from '$lib/api/instance/v1/instance_pb';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import * as Sheet from '$lib/components/ui/sheet';
	import { m } from '$lib/paraglide/messages';

	import { DataTable } from './restore/data-table';
</script>

<script lang="ts">
	let {
		virtualMachine,
		scope,
		reloadManager
	}: { virtualMachine: VirtualMachine; scope: string; reloadManager: ReloadManager } = $props();
</script>

<div class="flex items-center justify-end gap-1">
	<Sheet.Root>
		<Sheet.Trigger class="flex items-center gap-1">
			<Icon icon="ph:arrow-counter-clockwise" />
			{m.restore()}
		</Sheet.Trigger>

		<Sheet.Content class="min-w-[70vw] p-4">
			<DataTable {virtualMachine} {scope} {reloadManager} />
		</Sheet.Content>
	</Sheet.Root>
</div>
