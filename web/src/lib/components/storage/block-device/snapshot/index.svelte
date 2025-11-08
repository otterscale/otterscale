<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { setContext } from 'svelte';

	import { DataTable } from './data-table';

	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import * as Sheet from '$lib/components/ui/sheet';
</script>

<script lang="ts">
	let {
		image
	}: {
		image: Image;
	} = $props();

	setContext('image', image);
</script>

<div class="flex items-center justify-end gap-1">
	{image.snapshots.length}
	<Sheet.Root>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[38vw] p-4">
			<DataTable {image} />
		</Sheet.Content>
	</Sheet.Root>
</div>
