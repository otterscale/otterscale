<script lang="ts">
	import Icon from '@iconify/svelte';

	import { buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';

	import { listInstances, type pbInstance } from '$lib/pb';
	import {
		CreateOverview,
		CreateApplication,
		CreateCluster,
		CreateMachine
	} from '$lib/components/otterscale';

	import { onMount } from 'svelte';
	import { cn } from '$lib/utils';
	import { machines } from './create';

	let items = $state([
		{
			name: 'Machine',
			icon: 'ph:atom',
			active: false
		},
		{
			name: 'Cluster',
			icon: 'ph:graph',
			active: false
		},
		{
			name: 'Application',
			icon: 'ph:infinity',
			active: false
		}
	]);

	let open = $state(false);

	let pbMachines: pbInstance[] = $state([]);
	let pbClusters: pbInstance[] = $state([]);
	let pbApplications: pbInstance[] = $state([]);

	onMount(async () => {
		pbMachines = await listInstances(`type='MAAS'`);
		pbClusters = await listInstances(`type='JUJU'`);
		pbApplications = await listInstances(`type='kubernetes'`);
	});
	//
</script>

<Dialog.Root
	bind:open
	onOpenChange={(open) => {
		if (!open) {
			setTimeout(() => {
				items = items.map((item) => ({ ...item, active: false }));
			}, 100);
		}
	}}
>
	<Dialog.Trigger class={cn(buttonVariants({ size: 'icon' }), '[&_svg]:size-5')}>
		<Icon icon="ph:plus" />
	</Dialog.Trigger>
	<Dialog.Content class="max-h-[100vh] max-w-2xl overflow-y-auto">
		<Dialog.Header class="flex-col space-y-8 py-4">
			<Dialog.Title class="flex">
				{#each items as item}
					{#if item.active}
						<div class="flex items-center pl-2">
							<Icon icon={item.icon} class="size-8" />
							<div class="space-x-2 pl-2">{item.name}</div>
						</div>
					{/if}
				{/each}
				{#if items.filter((item) => item.active).length === 0}
					<div class="flex w-full items-center justify-center">Create Instance</div>
				{/if}
			</Dialog.Title>
			<Dialog.Description class="flex w-full justify-center px-2">
				{#if items[0].active}
					<CreateMachine bind:parent={open} items={machines} />
				{:else if items[1].active}
					<CreateCluster bind:parent={open} items={machines} />
				{:else if items[2].active}
					<CreateApplication bind:parent={open} items={machines} />
				{:else}
					<CreateOverview bind:items />
				{/if}
			</Dialog.Description>
		</Dialog.Header>
	</Dialog.Content>
</Dialog.Root>
