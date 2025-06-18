<script lang="ts">
	import { DataTable } from './data-table';
	import * as Sheet from '$lib/components/ui/sheet';
	import Icon from '@iconify/svelte';
	import { cn } from '$lib/utils';
</script>

<Sheet.Root>
	<Sheet.Trigger class={cn('flex h-full w-full items-center gap-2')}>
		<Icon icon="ph:cube" />
		Subvolume Group
	</Sheet.Trigger>
	<Sheet.Content class="flex min-w-[77vw] flex-col justify-center p-4">
		<DataTable />
	</Sheet.Content>
</Sheet.Root>
