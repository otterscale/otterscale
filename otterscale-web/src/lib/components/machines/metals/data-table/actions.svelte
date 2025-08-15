<script lang="ts" module>
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import { Layout } from '$lib/components/custom/data-table';
	import Add from './add.svelte';
	import PowerOff from './power-off.svelte';
	import Remove from './remove.svelte';
</script>

<script lang="ts">
	let {
		machine
	}: {
		machine: Machine;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Add {machine} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Remove {machine} />
	</Layout.ActionItem>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<PowerOff {machine} />
	</Layout.ActionItem>
</Layout.Actions>
