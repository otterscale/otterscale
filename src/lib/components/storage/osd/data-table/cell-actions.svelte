<script lang="ts" module>
	import type { ObjectStorageDaemon } from '$lib/api/storage/v1/storage_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import { m } from '$lib/paraglide/messages';

	import DoSMART from './action-do-SMART.svelte';
</script>

<script lang="ts">
	let { osd, scope }: { osd: ObjectStorageDaemon; scope: string } = $props();

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Actions.List bind:open>
	<Actions.Label>
		{m.actions()}
	</Actions.Label>
	<Actions.Separator />
	<Actions.Item>
		<DoSMART {osd} {scope} closeActions={close} />
	</Actions.Item>
</Actions.List>
