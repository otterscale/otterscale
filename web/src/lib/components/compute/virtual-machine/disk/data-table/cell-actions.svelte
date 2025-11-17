<script lang="ts" module>
	import type { EnhancedDisk } from '$lib/components/compute/virtual-machine/units/type';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Detach from './action-detach.svelte';
</script>

<script lang="ts">
	let {
		enhancedDisk,
		scope,
		reloadManager
	}: { enhancedDisk: EnhancedDisk; scope: string; reloadManager: ReloadManager } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Detach {enhancedDisk} {scope} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
