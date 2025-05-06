<script lang="ts">
	import { page } from '$app/state';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { capitalizeFirstLetter } from 'better-auth';
	import { toast } from 'svelte-sonner';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { Badge } from '$lib/components/ui/badge';
	import * as Table from '$lib/components/ui/table';
	import * as Tabs from '$lib/components/ui/tabs';
	import { Nexus, type Machine, type Network } from '$gen/api/nexus/v1/nexus_pb';
	import { formatCapacity } from '$lib/formatter';
	import { getContext, onMount, onDestroy } from 'svelte';
	import { writable } from 'svelte/store';
	import * as Alert from '$lib/components/ui/alert/index.js';

	const nodeType = 'MAAS';

	let {
		machine
	}: {
		machine: Machine;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const machineStore = writable<Machine>();
	async function refreshMachine() {
		while (page.url.searchParams.get('intervals')) {
			await new Promise((resolve) =>
				setTimeout(resolve, 1000 * Number(page.url.searchParams.get('intervals')))
			);
			console.log(`Refresh machine ${page.params.system_id}`);

			try {
				const response = await client.getMachine({ id: machine.id });
				machineStore.set(response);

				machine = $machineStore;
			} catch (error) {
				console.error('Error fetching machine:', error);
			}
		}
	}

	onMount(async () => {
		try {
			refreshMachine();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<div class="grid gap-3 space-y-3">
	{#if machine.statusMessage !== 'Deployed'}
		<Alert.Root>
			<Icon icon="ph:spinner" class="size-4 animate-spin" />
			<Alert.Title>{machine.status}</Alert.Title>
			<Alert.Description>{machine.statusMessage}</Alert.Description>
		</Alert.Root>
	{/if}

	{@render Summary()}

	<Tabs.Root value="hardware_information">
		<div class="flex justify-between">
			<Tabs.List class="w-fit rounded-sm">
				<Tabs.Trigger value="hardware_information">Hardware Information</Tabs.Trigger>
				<Tabs.Trigger value="block_device">Block Devices</Tabs.Trigger>
				<Tabs.Trigger value="network">Networks</Tabs.Trigger>
			</Tabs.List>
			{#if machine.workloadAnnotations && machine.workloadAnnotations['juju-model-uuid']}
				<Button href="/management/facility?scope={machine.workloadAnnotations['juju-model-uuid']}">
					<Icon icon="ph:rocket-launch" />
					Facility
				</Button>
			{/if}
		</div>
		<div>
			<Tabs.Content value="hardware_information" class="p-2">
				{@render TabContent_HardwareInformation()}
			</Tabs.Content>
			<Tabs.Content value="block_device">
				{@render TabContent_BlockDevices()}
			</Tabs.Content>
			<Tabs.Content value="network">
				{@render TabContent_Networks()}
			</Tabs.Content>
		</div>
	</Tabs.Root>
</div>

{#snippet Summary()}
	{@const formattedMemory = formatCapacity(machine.memoryMb)}
	{@const formattedStorage = formatCapacity(machine.storageMb)}
	<div class="grid grid-cols-5 gap-3">
		<Card.Root>
			<Card.Header class="h-10">
				<Card.Title>
					<div class="flex justify-between text-xs">
						MACHINE
						<div class="font-light text-muted-foreground">
							{machine.id}
						</div>
					</div>
				</Card.Title>
			</Card.Header>
			<Card.Content class="h-20">
				<div class="text-2xl">
					{machine.fqdn}
				</div>
			</Card.Content>
			<Card.Footer>
				<div class="flex flex-wrap gap-1">
					{#if machine.powerType !== ''}
						<Badge variant="outline">{machine.powerType}</Badge>
					{/if}
					{#each machine.tags as tag}
						<Badge variant="outline" class="text-muted-foreground">{tag}</Badge>
					{/each}
				</div>
			</Card.Footer>
		</Card.Root>
		<Card.Root>
			<Card.Header class="h-10">
				<Card.Title
					><div class="flex justify-between text-xs">
						STATUS
						{#if machine.powerState === 'on'}
							<Badge>Power On</Badge>
						{:else}
							<Badge variant="destructive">Power Off</Badge>
						{/if}
					</div></Card.Title
				>
			</Card.Header>
			<Card.Content class="h-20">
				<div class="text-2xl">
					{machine.status}
				</div>
			</Card.Content>
			<Card.Footer>
				<div class="truncate text-xs font-light text-muted-foreground">
					{capitalizeFirstLetter(machine.osystem)}
					{machine.hweKernel}
					{capitalizeFirstLetter(machine.distroSeries)}
				</div>
			</Card.Footer>
		</Card.Root>
		<Card.Root>
			<Card.Header class="h-10">
				<Card.Title
					><div class="flex justify-between text-xs">
						CPU
						<div class="font-light text-muted-foreground">
							{machine.architecture}
						</div>
					</div>
				</Card.Title>
			</Card.Header>
			<Card.Content class="h-20">
				<div class="text-2xl">
					{machine.cpuCount} cores
				</div>
			</Card.Content>
			<Card.Footer>
				<div class="truncate text-xs font-light text-muted-foreground">
					{machine.hardwareInformation.cpu_model}
				</div>
			</Card.Footer>
		</Card.Root>
		<Card.Root>
			<Card.Header class="h-10">
				<Card.Title><div class="flex justify-between text-xs">MEMORY</div></Card.Title>
			</Card.Header>
			<Card.Content class="flex h-20 items-end gap-1">
				<p class="text-4xl">
					{formattedMemory.value}
				</p>
				<p class="text-2xl">
					{formattedMemory.unit}
				</p>
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header class="h-10">
				<Card.Title><div class="flex justify-between text-xs">STORAGE</div></Card.Title>
			</Card.Header>
			<Card.Content class="h-20 space-y-1">
				<div class="flex items-end gap-1">
					<p class="text-4xl">
						{formattedStorage.value}
					</p>
					<p class="text-2xl">
						{formattedStorage.unit}
					</p>
				</div>
				<p class="truncate text-xs font-light text-muted-foreground">
					over {machine.blockDevices.length} disks
				</p>
			</Card.Content>
		</Card.Root>
	</div>
{/snippet}

{#snippet TabContent_HardwareInformation()}
	<div class="grid gap-4">
		<div class="grid gap-2">
			<div class="flex items-center gap-1">
				<Icon icon="ph:desktop" class="size-5" />
				<Label class="text-base">System</Label>
			</div>
			<Table.Root>
				<Table.Header class="bg-muted/50">
					<Table.Row class="*:text-xs *:font-light">
						<Table.Head>VENDOR</Table.Head>
						<Table.Head>PRODUCT</Table.Head>
						<Table.Head>VERSION</Table.Head>
						<Table.Head>SERIAL</Table.Head>
						<Table.Head>SKU</Table.Head>
						<Table.Head>FAMILY</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					<Table.Row class="*:text-xs">
						<Table.Cell>{machine.hardwareInformation.system_vendor}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.system_product}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.system_version}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.system_serial}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.system_sku}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.system_family}</Table.Cell>
					</Table.Row>
				</Table.Body>
			</Table.Root>
		</div>
		<div class="grid gap-2">
			<div class="flex items-center gap-1">
				<Icon icon="ph:circuitry" class="size-5" />
				<Label class="text-base">Mainboard</Label>
			</div>

			<Table.Root>
				<Table.Header class="bg-muted/50">
					<Table.Row class="*:text-xs *:font-light">
						<Table.Head>VENDOR</Table.Head>
						<Table.Head>PRODUCT</Table.Head>
						<Table.Head>FIRMWARE</Table.Head>
						<Table.Head>BOOT MODE</Table.Head>
						<Table.Head>VERSION</Table.Head>
						<Table.Head>DATE</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					<Table.Row class="*:text-xs">
						<Table.Cell>{machine.hardwareInformation.mainboard_vendor}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.mainboard_product}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.mainboard_firmware_vendor}</Table.Cell>
						<Table.Cell>{machine.biosBootMethod.toUpperCase()}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.mainboard_firmware_version}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.mainboard_firmware_date}</Table.Cell>
					</Table.Row>
				</Table.Body>
			</Table.Root>
		</div>
		<div class="grid gap-2">
			<div class="flex items-center gap-1">
				<Icon icon="ph:computer-tower" class="size-5" />
				<Label class="text-base">Chassis</Label>
			</div>

			<Table.Root>
				<Table.Header class="bg-muted/50">
					<Table.Row class="*:text-xs *:font-light">
						<Table.Head>VENDOR</Table.Head>
						<Table.Head>TYPE</Table.Head>
						<Table.Head>VERSION</Table.Head>
						<Table.Head>SERIAL</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					<Table.Row class="*:text-xs">
						<Table.Cell>{machine.hardwareInformation.chassis_vendor}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.chassis_type}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.chassis_version}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.chassis_serial}</Table.Cell>
					</Table.Row>
				</Table.Body>
			</Table.Root>
		</div>
	</div>
{/snippet}
{#snippet TabContent_BlockDevices()}
	<Table.Root>
		<Table.Header>
			<Table.Row class="*:text-xs *:font-light">
				<Table.Head>NAME</Table.Head>
				<Table.Head>MODEL</Table.Head>
				<Table.Head>SERIAL</Table.Head>
				<Table.Head>BOOT DISK</Table.Head>
				<Table.Head>FIRMWARE VERSION</Table.Head>
				<Table.Head>TYPE</Table.Head>
				<Table.Head>USER FOR</Table.Head>
				<Table.Head>TAGS</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each machine.blockDevices as blockDevice}
				<Table.Row class="*:text-xs">
					<Table.Cell>{blockDevice.name}</Table.Cell>
					<Table.Cell>{blockDevice.model}</Table.Cell>
					<Table.Cell>{blockDevice.serial}</Table.Cell>
					<Table.Cell>{blockDevice.bootDisk}</Table.Cell>
					<Table.Cell>{blockDevice.firmwareVersion}</Table.Cell>
					<Table.Cell>{blockDevice.type}</Table.Cell>
					<Table.Cell>{blockDevice.usedFor}</Table.Cell>
					<Table.Cell>
						<div class="flex space-x-1">
							{#each blockDevice.tags as tag}
								<Badge variant="outline">
									{tag}
								</Badge>
							{/each}
						</div>
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
{/snippet}
{#snippet TabContent_Networks()}
	<Table.Root>
		<Table.Header>
			<Table.Row class="*:text-xs *:font-light [&>th]:py-2 [&>th]:align-top">
				<Table.Head>
					NAME
					<p class="text-muted-foreground">MAC Address</p>
				</Table.Head>
				<Table.Head>
					IP ADDRESS
					<p class="text-muted-foreground">Subnet</p>
				</Table.Head>
				<Table.Head>
					LINK SPEED
					<p class="text-muted-foreground">Link Connected</p>
				</Table.Head>
				<Table.Head>
					FABRIC
					<p class="text-muted-foreground">VLAN</p>
				</Table.Head>
				<Table.Head>TYPE</Table.Head>
				<Table.Head>DHCP ON</Table.Head>
				<Table.Head>BOOT INTERFACE</Table.Head>
				<Table.Head>INTERFACE SPEED</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each machine.networkInterfaces as networkInterface}
				<Table.Row class="*:text-xs [&>td]:align-top">
					<Table.Cell>
						<p>{networkInterface.name}</p>
						<p>{networkInterface.macAddress}</p>
					</Table.Cell>
					<Table.Cell>
						<p>
							{networkInterface.ipAddress}
						</p>
						<p>
							{networkInterface.subnetName}
						</p>
					</Table.Cell>
					<Table.Cell>
						<p>{networkInterface.linkSpeed} Mbps</p>
						<Icon
							icon={networkInterface.linkConnected ? 'ph:check-circle' : 'ph:x-circle'}
							style="color: {networkInterface.linkConnected ? 'green' : 'red'}"
						/>
					</Table.Cell>
					<Table.Cell>
						<p>{networkInterface.fabricName}</p>
						<p>
							{networkInterface.vlanName}
						</p>
					</Table.Cell>
					<Table.Cell>{networkInterface.type}</Table.Cell>
					<Table.Cell class="align-center">
						<Icon
							icon={networkInterface.dhcpOn ? 'ph:check' : 'ph:x'}
							style="color: {networkInterface.dhcpOn ? 'green' : 'red'}"
						/>
					</Table.Cell>
					<Table.Cell>{networkInterface.bootInterface}</Table.Cell>
					<Table.Cell>{networkInterface.interfaceSpeed} Mbps</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
{/snippet}
