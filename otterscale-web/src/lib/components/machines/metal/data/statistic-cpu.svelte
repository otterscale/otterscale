<script lang="ts" module>
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';
	import { Statistic } from '../layout';
</script>

<script lang="ts">
	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<Statistic.Root>
	<Statistic.Header>
		<Statistic.Title>CPU</Statistic.Title>
		<Statistic.Action>
			<Icon icon="ph:cpu" />
			{$machine.architecture}
		</Statistic.Action>
	</Statistic.Header>
	<Statistic.Content>
		{$machine.cpuCount} cores
	</Statistic.Content>
	<Statistic.Footer>
		{$machine.hardwareInformation.cpu_model}
	</Statistic.Footer>
</Statistic.Root>
