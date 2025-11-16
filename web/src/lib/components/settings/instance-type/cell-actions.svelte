<script lang="ts" module>
	import type { InstanceType } from '$lib/api/instance/v1/instance_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
</script>

<script lang="ts">
	let { instanceType, scope, reloadManager }: { instanceType: InstanceType; scope: string; reloadManager: ReloadManager } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Delete {instanceType} {scope} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
