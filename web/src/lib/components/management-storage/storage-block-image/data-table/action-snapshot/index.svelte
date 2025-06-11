<script lang="ts">
	import { DataTable } from './data-table';
	import Create from './create.svelte';
	import * as Drawer from '$lib/components/ui/drawer';
	import Icon from '@iconify/svelte';
	import { cn } from '$lib/utils';
</script>

<Drawer.Root direction="bottom">
	<Drawer.Trigger class={cn('flex h-full w-full items-center gap-2')}>
		<Icon icon="ph:camera" />
		Snapshot
	</Drawer.Trigger>
	<Drawer.Content class="max-h-[77vh] p-4">
		<Create />
		<div class="overflow-y-auto">
			<DataTable />
		</div>
	</Drawer.Content>
</Drawer.Root>
