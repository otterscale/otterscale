<script lang="ts" module>
	import DoSMART from './action-do-SMART.svelte';

	import type { OSD } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { osd }: { osd: OSD } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<DoSMART {osd} />
	</Layout.ActionItem>
</Layout.Actions>
