<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { DataTable } from './data-table';
	import type { BlockImage } from '../types';

	let { blockImage }: { blockImage: BlockImage } = $props();
</script>

<Sheet.Root>
	<Sheet.Trigger class={cn('flex h-full w-full items-center gap-2')}>
		<Icon icon="ph:camera" />
		Snapshot
	</Sheet.Trigger>
	<Sheet.Content class="flex min-w-[77vw] flex-col justify-center p-4">
		<DataTable {blockImage} />
	</Sheet.Content>
</Sheet.Root>
