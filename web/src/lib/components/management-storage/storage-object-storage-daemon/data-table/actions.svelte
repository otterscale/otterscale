<script lang="ts">
	import type { OSD } from '$gen/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import DoSMART from './doSMART.svelte';

	let { osd }: { osd: OSD } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<DoSMART {osd} />
	</Layout.ActionItem>
</Layout.Actions>
