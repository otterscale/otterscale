<script lang="ts" module>
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Layout } from '$lib/components/custom/instance';
	import { Badge } from '$lib/components/ui/badge';
	import Icon from '@iconify/svelte';
	import { capitalizeFirstLetter } from 'better-auth';
	import { type Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<Layout.Statistic.Root>
	<Layout.Statistic.Header>
		<Layout.Statistic.Title>
			<Badge variant={$machine.powerState === 'on' ? 'default' : 'destructive'}>
				<Icon icon="ph:power" />
				Power {$machine.powerState}
			</Badge>
		</Layout.Statistic.Title>
		<Layout.Statistic.Action>
			<Badge variant="outline">
				<Icon icon="ph:power" />
				{$machine.powerType}
			</Badge>
		</Layout.Statistic.Action>
	</Layout.Statistic.Header>
	<Layout.Statistic.Content>
		{$machine.status}
	</Layout.Statistic.Content>
	<Layout.Statistic.Footer>
		{capitalizeFirstLetter($machine.osystem)}
		{$machine.hweKernel}
		{capitalizeFirstLetter($machine.distroSeries)}
	</Layout.Statistic.Footer>
</Layout.Statistic.Root>
