<script lang="ts" module>
	import PowerOff from './action-power-off.svelte';
	import Remove from './action-remove.svelte';

	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		machine
	}: {
		machine: Machine;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem
		disabled={(machine.status.toLowerCase() !== 'ready' &&
			machine.status.toLowerCase() !== 'releasing') ||
			!!machine.workloadAnnotations['juju-is-controller']}
	>
		<Remove {machine} />
	</Layout.ActionItem>
	<Layout.ActionItem
		disabled={machine.powerState.toLowerCase() !== 'on' ||
			machine.status.toLowerCase() === 'commissioning' ||
			machine.status.toLowerCase() === 'testing' ||
			machine.status.toLowerCase() === 'deploying'}
	>
		<PowerOff {machine} />
	</Layout.ActionItem>
</Layout.Actions>
