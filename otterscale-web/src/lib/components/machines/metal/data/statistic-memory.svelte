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

	const { value, unit } = formatCapacity(Number($machine.memoryMb)*1000*1000);
</script>

<Layout.Statistic.Root>
	<Layout.Statistic.Header>
		<Layout.Statistic.Title>MEMORY</Layout.Statistic.Title>
	</Layout.Statistic.Header>
	<Layout.Statistic.Content>
		{value}
		{unit}
	</Layout.Statistic.Content>
	<Layout.Statistic.Footer></Layout.Statistic.Footer>
</Layout.Statistic.Root>
