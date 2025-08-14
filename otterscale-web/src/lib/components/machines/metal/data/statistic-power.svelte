<script lang="ts" module>
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Badge } from '$lib/components/ui/badge';
	import Icon from '@iconify/svelte';
	import { capitalizeFirstLetter } from 'better-auth';
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
		<Statistic.Title>
			<Badge variant={$machine.powerState === 'on' ? 'default' : 'destructive'}>
				<Icon icon="ph:power" />
				Power {$machine.powerState}
			</Badge>
		</Statistic.Title>
		<Statistic.Action>
			<Badge variant="outline">
				<Icon icon="ph:power" />
				{$machine.powerType}
			</Badge>
		</Statistic.Action>
	</Statistic.Header>
	<Statistic.Content>
		{$machine.status}
	</Statistic.Content>
	<Statistic.Footer>
		{capitalizeFirstLetter($machine.osystem)}
		{$machine.hweKernel}
		{capitalizeFirstLetter($machine.distroSeries)}
	</Statistic.Footer>
</Statistic.Root>
