<script lang="ts" module>
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Layout } from '$lib/components/custom/instance';
	import { m } from '$lib/paraglide/messages';
	import Icon from '@iconify/svelte';
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
		<Layout.Statistic.Title>{m.cpu()}</Layout.Statistic.Title>
		<Layout.Statistic.Action>
			<Icon icon="ph:cpu" />
			{$machine.architecture}
		</Layout.Statistic.Action>
	</Layout.Statistic.Header>
	<Layout.Statistic.Content>
		{$machine.cpuCount}
		{m.core()}
	</Layout.Statistic.Content>
	<Layout.Statistic.Footer>
		{$machine.hardwareInformation.cpu_model}
	</Layout.Statistic.Footer>
</Layout.Statistic.Root>
