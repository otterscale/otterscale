<script lang="ts" module>
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Layout } from '$lib/components/custom/instance';
	import { formatCapacity } from '$lib/formatter';
	import { type Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();

	const { value, unit } = formatCapacity($machine.memoryMb);
</script>

<Layout.Statistic.Root>
	<Layout.Statistic.Header>
		<Layout.Statistic.Title>STORAGE</Layout.Statistic.Title>
	</Layout.Statistic.Header>
	<Layout.Statistic.Content>
		{value}
		{unit}
	</Layout.Statistic.Content>
	<Layout.Statistic.Footer>
		over {$machine.blockDevices.length} disks
	</Layout.Statistic.Footer>
</Layout.Statistic.Root>
