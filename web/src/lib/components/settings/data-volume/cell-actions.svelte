<script lang="ts" module>
	import Delete from './action-delete.svelte';
	import Extend from './action-extend.svelte';

	import type { DataVolume } from '$lib/api/instance/v1/instance_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { dataVolume }: { dataVolume: DataVolume } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Extend {dataVolume} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {dataVolume} />
	</Layout.ActionItem>
</Layout.Actions>
