<script lang="ts">
	import type { OSD } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import DoSMART from './action-do-SMART.svelte';

	let { osd }: { osd: OSD } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<DoSMART {osd} />
	</Layout.ActionItem>
</Layout.Actions>
