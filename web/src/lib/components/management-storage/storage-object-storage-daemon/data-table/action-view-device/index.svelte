<script lang="ts">
	import { DataTable } from './data-table';
	import * as Sheet from '$lib/components/ui/sheet';
	import Icon from '@iconify/svelte';
	import { cn } from '$lib/utils';
	import type { ObjectStorageDaemon } from '../types';

	let { row }: { row: ObjectStorageDaemon } = $props();
</script>

<Sheet.Root>
	<Sheet.Trigger>
		<Icon icon="ph:arrow-square-out" />
	</Sheet.Trigger>
	<Sheet.Content class="flex min-w-[77vw] flex-col justify-center p-4">
		<DataTable />
	</Sheet.Content>
</Sheet.Root>
