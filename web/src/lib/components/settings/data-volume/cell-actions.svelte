<script lang="ts" module>
	import type { DataVolume } from '$lib/api/instance/v1/instance_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Extend from './action-extend.svelte';
</script>

<script lang="ts">
	let { dataVolume, scope, reloadManager }: { dataVolume: DataVolume; scope: string; reloadManager: ReloadManager  } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Extend {dataVolume} {scope} {reloadManager} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {dataVolume} {scope} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
