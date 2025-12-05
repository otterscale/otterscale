<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';

	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Layout } from '$lib/components/custom/instance';
	import { Badge } from '$lib/components/ui/badge';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<Layout.Statistic.Root class="relative overflow-hidden">
	{#if ['Commissioning', 'Deploying', 'Releasing'].includes($machine.status)}
		<Spinner class="absolute right-4 bottom-4 size-8 text-primary/20" />
	{/if}
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
		<div class="flex flex-col">
			<span class=" flex items-center gap-1">
				{$machine.status}
			</span>
			{#if (['Commissioning', 'Deploying', 'Releasing'].includes($machine.status) || $machine.status.startsWith('Failed')) && $machine.statusMessage !== $machine.status}
				<p class="mt-1 h-0 text-sm text-muted-foreground">
					<span class="flex items-center gap-1">
						{$machine.statusMessage}
					</span>
				</p>
			{/if}
		</div>
	</Layout.Statistic.Content>
	<Layout.Statistic.Footer class="gap-1">
		{#if $machine.osystem}
			<span class="capitalize">{$machine.osystem}</span>
		{/if}
		{#if $machine.hweKernel}
			<span class="uppercase">{$machine.hweKernel}</span>
		{/if}
		{#if $machine.distroSeries}
			<span class="capitalize">{$machine.distroSeries}</span>
		{/if}
	</Layout.Statistic.Footer>
</Layout.Statistic.Root>
