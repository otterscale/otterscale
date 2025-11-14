<script lang="ts" module>
	import type { OSD } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	import DoSMART from './action-do-SMART.svelte';
</script>

<script lang="ts">
	let { osd, scope }: { osd: OSD; scope: string } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<DoSMART {osd} {scope} />
	</Layout.ActionItem>
</Layout.Actions>
