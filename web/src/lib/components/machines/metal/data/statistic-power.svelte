<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';

	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Layout } from '$lib/components/custom/instance';
	import { Badge } from '$lib/components/ui/badge';
	import { m } from '$lib/paraglide/messages';
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
			{$machine.powerType}
		</Layout.Statistic.Title>
		<Layout.Statistic.Action>
			<Badge variant={$machine.powerState === 'on' ? 'default' : 'destructive'}>
				<Icon icon="ph:power" />
				{m.power()}
				{$machine.powerState}
			</Badge>
		</Layout.Statistic.Action>
	</Layout.Statistic.Header>
	<Layout.Statistic.Content>
		{$machine.status}
	</Layout.Statistic.Content>
	<Layout.Statistic.Footer>
		<sapn class="capitalize"> {$machine.osystem}</sapn>
		{$machine.hweKernel}
		<sapn class="capitalize"> {$machine.distroSeries}</sapn>
	</Layout.Statistic.Footer>
</Layout.Statistic.Root>
