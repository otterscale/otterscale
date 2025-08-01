<script lang="ts" module>
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { type Writable } from 'svelte/store';
	import { Statistic } from '../layout';
	import { formatCapacity } from '$lib/formatter';
</script>

<script lang="ts">
	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();

	const { value, unit } = formatCapacity($machine.memoryMb);
</script>

<Statistic.Root>
	<Statistic.Header>
		<Statistic.Title>STORAGE</Statistic.Title>
	</Statistic.Header>
	<Statistic.Content>
		{value}
		{unit}
	</Statistic.Content>
	<Statistic.Footer>
		over {$machine.blockDevices.length} disks
	</Statistic.Footer>
</Statistic.Root>
