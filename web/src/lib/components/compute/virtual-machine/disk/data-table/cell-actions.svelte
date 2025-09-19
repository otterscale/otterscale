<script lang="ts" module>
	import Extend from './action-extend.svelte';

	import type { VirtualMachineDisk } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { virtualMachineDisk }: { virtualMachineDisk: VirtualMachineDisk } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Extend {virtualMachineDisk} />
	</Layout.ActionItem>
</Layout.Actions>
