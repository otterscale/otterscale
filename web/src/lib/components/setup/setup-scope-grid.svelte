<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { resolve } from '$app/paths';
	import {
		type Facility,
		type Facility_Status,
		type Facility_Unit,
		FacilityService
	} from '$lib/api/facility/v1/facility_pb';
	import * as Accordion from '$lib/components/ui/accordion';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { Separator } from '$lib/components/ui/separator';
	import { Switch } from '$lib/components/ui/switch';
	import { m } from '$lib/paraglide/messages';

	import SetupNodeGPUMode from './setup-node-gpu-mode.svelte';

	let {
		services,
		scope,
		facilities,
		autoRefresh = $bindable(true)
	}: {
		services: Record<
			string,
			{
				name: string;
				icon: string;
				title: string;
				gridClass: string;
			}
		>;
		scope: string;
		facilities: Writable<Facility[]>;
		autoRefresh: boolean;
	} = $props();

	const transport: Transport = getContext('transport');
	const facilityClient = createClient(FacilityService, transport);

	// Helper functions
	function findFacilityByService(serviceName: string): Facility | undefined {
		return $facilities.find(
			(facility) => facility.name.includes(serviceName) && facility.units.length > 0
		);
	}

	function countUnitsByService(serviceName: string): number {
		const facility = findFacilityByService(serviceName);
		return facility?.units.filter((u) => u.workloadStatus?.state === 'active').length ?? 0;
	}

	function getStatusClass(status: Facility_Status | undefined): string {
		return status?.state === 'active' ? 'text-muted-foreground' : 'font-semibold text-red-500';
	}
</script>

<div class="grid w-full grid-cols-3 gap-4 sm:gap-6 lg:grid-cols-6">
	<div class="col-span-3 flex justify-end space-x-4 rounded-lg sm:space-x-6 lg:col-span-6">
		<Button variant="ghost" disabled>
			<Icon icon="ph:plus" class="size-4" />
			{m.add_node()}
		</Button>
		<div class="flex items-center space-x-2">
			<Switch id="auto-update" bind:checked={autoRefresh} />
			<Label for="auto-update">{m.auto_update()}</Label>
		</div>
	</div>

	{#each Object.values(services) as service}
		{@const facility = findFacilityByService(service.name)}
		{@const count = countUnitsByService(service.name)}

		<div
			class="relative bg-muted {service.gridClass} flex flex-col space-y-2 overflow-hidden rounded-lg p-4 shadow-sm md:p-6 lg:p-10"
		>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={service.icon} class="size-72" />
			</div>

			<div class="z-10 flex items-center justify-between">
				<div class="flex flex-col items-start">
					<h2 class="text-2xl font-medium">{service.title}</h2>
					<p class="text-sm capitalize {getStatusClass(facility?.status)}">
						{facility?.status?.details}
					</p>
				</div>

				<div class="flex space-x-1 text-3xl sm:mb-2 lg:text-5xl">
					<span>{count}</span>
					{#if facility && facility.units.length > count}
						<Icon icon="ph:arrow-up-bold" class="size-6 animate-bounce" />
					{/if}
				</div>
			</div>

			<Separator class="z-10" />

			{#if facility}
				<div class="z-10 pb-8 sm:pb-12">
					{@render facilityDisplay(facility)}
				</div>
			{/if}
		</div>
	{/each}
</div>

{#snippet facilityDisplay(facility: Facility)}
	<Accordion.Root
		type="multiple"
		value={facility.units
			.filter((unit) => unit.workloadStatus?.state !== 'active')
			.map((unit) => unit.name)}
	>
		{#each facility.units.sort((a, b) => a.name.localeCompare(b.name)) as unit}
			<Accordion.Item value={unit.name}>
				<Accordion.Trigger class="py-4">
					<div class="flex flex-col space-y-1">
						<div class="flex items-center space-x-2">
							<span class="text-xs font-medium md:text-base lg:text-lg">
								{unit.name}
							</span>

							<button
								class="hover:cursor-pointer"
								onclick={(e) => {
									e.stopPropagation();
									toast.promise(
										() =>
											facilityClient.resolveFacilityUnitErrors({
												scope: scope,
												unitName: unit.name
											}),
										{
											loading: 'Loading...',
											success: () => {
												facilityClient
													.listFacilities({
														scope: scope
													})
													.then((response) => {
														facilities.set(response.facilities);
													});

												return `Resolve ${unit.name} success`;
											},
											error: (e) => {
												let msg = `Fail to resolve ${unit.name}`;
												toast.error(msg, {
													description: (e as ConnectError).message.toString(),
													duration: Number.POSITIVE_INFINITY,
													closeButton: true
												});
												return msg;
											}
										}
									);
								}}
							>
								<Icon icon="ph:arrow-counter-clockwise" class="size-4" />
							</button>

							{#if unit.machineId}
								<a
									href={resolve('/(auth)/machines/metal/[id]', {
										id: unit.machineId
									})}
								>
									<Icon icon="ph:computer-tower" class="size-4" />
								</a>
							{/if}

							{#if unit.machineId}
								<SetupNodeGPUMode {unit} {scope} class="hover:cursor-pointer" />
							{/if}

							{#if unit.leader}
								<Icon icon="ph:star-fill" class="size-4 text-yellow-400" />
							{/if}
						</div>

						<div class="flex items-center space-x-2 font-normal">
							<span class="text-sm leading-none text-muted-foreground">
								{unit.version !== '' ? unit.version : '-'}
							</span>

							<span class="text-sm leading-none {getStatusClass(unit.workloadStatus)}">
								{unit.workloadStatus?.details}
							</span>
						</div>
					</div>
				</Accordion.Trigger>
				<Accordion.Content class="space-y-4">
					{@render subordinatesDisplay(
						unit.subordinates.sort((a, b) => a.name.localeCompare(b.name))
					)}
				</Accordion.Content>
			</Accordion.Item>
		{/each}
	</Accordion.Root>
{/snippet}

{#snippet subordinatesDisplay(subordinates: Facility_Unit[])}
	<div class="flex flex-col space-y-2">
		{#each subordinates as subordinate}
			{#if subordinate.workloadStatus}
				<div
					class="flex items-center space-x-1 text-sm {getStatusClass(subordinate.workloadStatus)}"
				>
					<div class="truncate">
						[{subordinate.name}] -
						<span class="capitalize">{subordinate.workloadStatus.details}</span>
					</div>
				</div>
			{/if}
		{/each}
	</div>
{/snippet}
