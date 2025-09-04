<script lang="ts" module>
	import Add from './action-add.svelte';
	import PowerOff from './action-power-off.svelte';
	import Remove from './action-remove.svelte';

	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		machine,
	}: {
		machine: Machine;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem disabled={!!machine.workloadAnnotations['juju-model-uuid']}>
		<Add {machine} />
	</Layout.ActionItem>
	<Layout.ActionItem
		disabled={!!machine.workloadAnnotations['juju-is-controller'] ||
			!machine.workloadAnnotations['juju-model-uuid']}
	>
		<Remove {machine} />
	</Layout.ActionItem>
	<Layout.ActionSeparator />
	<Layout.ActionItem disabled={machine.powerState.toLowerCase() !== 'on'}>
		<PowerOff {machine} />
	</Layout.ActionItem>
</Layout.Actions>
