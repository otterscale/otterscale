<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { VirtualMachine } from '$lib/api/instance/v1/instance_pb';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as Tooltip from '$lib/components/ui/tooltip';
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

<div class="flex w-full items-center gap-1">
	<Sheet.Root>
		<!-- TODO: disabled until feature is implemented -->
		<!-- <Sheet.Trigger class="flex items-center gap-1">
			<Icon icon="ph:arrow-counter-clockwise" />
			{m.restore()}
		</Sheet.Trigger> -->
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger class="w-full">
					<Modal.Trigger variant="creative" disabled>
						<Icon icon="ph:arrow-counter-clockwise" />
						{m.restore()}
					</Modal.Trigger>
				</Tooltip.Trigger>
				<Tooltip.Content>{m.under_development()}</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>

		<Sheet.Content class="min-w-[70vw] p-4">
			<DataTable {virtualMachine} {scope} {reloadManager} />
		</Sheet.Content>
	</Sheet.Root>
</div>
