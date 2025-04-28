<script lang="ts">
	// External dependencies
	import Icon from '@iconify/svelte';
	import { capitalizeFirstLetter } from 'better-auth';
	import { toast } from 'svelte-sonner';
	import * as Card from '$lib/components/ui/card/index.js';

	// Internal UI components
	import { Badge } from '$lib/components/ui/badge';
	import * as Table from '$lib/components/ui/table';
	import * as Tabs from '$lib/components/ui/tabs';

	// Internal utilities and types
	import { type Machine, type Network } from '$gen/api/nexus/v1/nexus_pb';
	import { formatCapacity } from '$lib/formatter';
	import { nodeIcon } from '$lib/node';

	const nodeType = 'MAAS';

	let {
		machine
	}: {
		machine: Machine;
	} = $props();
</script>

<div class="grid gap-3 space-y-3 p-3">
	{@render Identifier()}
	{@render Summary()}

	<Tabs.Root value="workload_annotation">
		<Tabs.List class={`grid w-fit grid-cols-4 rounded-sm`}>
			<Tabs.Trigger value="workload_annotation">Workload Annotation</Tabs.Trigger>
			<Tabs.Trigger value="hardware_information">Hardware Information</Tabs.Trigger>
			<Tabs.Trigger value="block_device">Block Devices</Tabs.Trigger>
			<Tabs.Trigger value="network">Networks</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="workload_annotation" class="h-full">
			{@render TabContent_WorkloadAnnotation()}
		</Tabs.Content>
		<Tabs.Content value="hardware_information">
			{@render TabContent_HardwareInformation()}
		</Tabs.Content>
		<Tabs.Content value="block_device">
			{@render TabContent_BlockDevices()}
		</Tabs.Content>
		<Tabs.Content value="network">
			{@render TabContent_Networks()}
		</Tabs.Content>
	</Tabs.Root>
</div>

{#snippet Identifier()}
	<div class="flex items-center">
		<div class="items-between flex space-x-1">
			<Icon icon={nodeIcon(nodeType)} class="h-full w-36" />
			<div class="flex flex-col justify-between">
				<div class="flex flex-col p-1">
					<div class="font-base text-xl">{machine.fqdn}</div>
					<div class="flex text-base text-muted-foreground">
						{machine.id}
					</div>
				</div>
				<p>{machine.description}</p>
				<div class="flex flex-wrap gap-1">
					{#each machine.tags as tag}
						<Badge variant="outline" class="text-muted-foreground">
							<Icon icon="ph:tag" class="size-5 sm:flex" />
							<span class="pl-1 text-sm">{tag}</span>
						</Badge>
					{/each}
				</div>
			</div>
		</div>
	</div>
{/snippet}

{#snippet Summary()}
	{@const formattedMemory = formatCapacity(machine.memoryMb)}
	{@const formattedStorage = formatCapacity(machine.storageMb / 1024)}
	<div class="grid grid-cols-5 gap-3 *:border-none *:shadow-none">
		<Card.Root>
			<Card.Header>
				<Card.Title><div class="flex justify-between text-xs">POWER</div></Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl">
					{capitalizeFirstLetter(machine.powerState)}
				</div>
			</Card.Content>
			<Card.Footer>
				<div class="truncate text-xs font-light text-muted-foreground">
					{machine.powerType}
				</div>
			</Card.Footer>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title><div class="flex justify-between text-xs">MACHINE STATUS</div></Card.Title>
			</Card.Header>
			<Card.Content>
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
			<Card.Header>
				<Card.Title
					><div class="flex justify-between text-xs">
						CPU
						<div class="font-light text-muted-foreground">
							{machine.architecture}
						</div>
					</div>
				</Card.Title>
			</Card.Header>
			<Card.Content>
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
			<Card.Header>
				<Card.Title><div class="flex justify-between text-xs">MEMORY</div></Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl">
					{formattedMemory.value}
					{formattedMemory.unit}
				</div>
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title><div class="flex justify-between text-xs">STORAGE</div></Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl">
					{formattedStorage.value}
					{formattedStorage.unit}
				</div>
			</Card.Content>
			<Card.Footer>
				<div class="truncate text-xs font-light text-muted-foreground">
					over {machine.blockDevices.length} disks
				</div>
			</Card.Footer>
		</Card.Root>
	</div>
{/snippet}

{#snippet TabContent_WorkloadAnnotation()}
	{#if Object.keys(machine.workloadAnnotations).length == 0}
		<p class="flex h-full flex-col items-center justify-center text-muted-foreground">
			There is no workload annotations.
		</p>
	{:else}
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>Juju Controller UUID</Table.Head>
					<Table.Head>Juju Machine ID</Table.Head>
					<Table.Head>Juju Model UUID</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				<Table.Row>
					<Table.Cell>
						{#if machine.workloadAnnotations['juju-controller-uuid']}
							<Badge variant="outline">
								{machine.workloadAnnotations['juju-controller-uuid']}
							</Badge>
						{/if}
					</Table.Cell>
					<Table.Cell>
						{#if machine.workloadAnnotations['juju-machine-id']}
							{machine.workloadAnnotations['juju-machine-id']}
						{/if}
					</Table.Cell>
					<Table.Cell>
						{#if machine.workloadAnnotations['juju-model-uuid']}
							<span class="flex items-center gap-1">
								<a href={`/management/scope/${machine.workloadAnnotations['juju-model-uuid']}`}>
									<Icon icon="ph:arrow-square-out" />
								</a>
								<Badge variant="outline">
									{machine.workloadAnnotations['juju-model-uuid']}
								</Badge>
							</span>
						{/if}
					</Table.Cell>
				</Table.Row>
			</Table.Body>
		</Table.Root>
	{/if}
{/snippet}
{#snippet TabContent_HardwareInformation()}
	<div class="flex flex-col gap-1 [&>fieldset]:rounded-lg [&>fieldset]:border [&>fieldset]:p-3">
		<fieldset>
			<legend>System</legend>
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>Vendor</Table.Head>
						<Table.Head>Product</Table.Head>
						<Table.Head>Version</Table.Head>
						<Table.Head>Serial</Table.Head>
						<Table.Head>SKU</Table.Head>
						<Table.Head>Family</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					<Table.Row>
						<Table.Cell>{machine.hardwareInformation.system_vendor}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.system_product}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.system_version}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.system_serial}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.system_sku}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.system_family}</Table.Cell>
					</Table.Row>
				</Table.Body>
			</Table.Root>
		</fieldset>
		<fieldset>
			<legend>Mainboard</legend>
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>Vendor</Table.Head>
						<Table.Head>Product</Table.Head>
						<Table.Head>Firmware</Table.Head>
						<Table.Head>Boot Mode</Table.Head>
						<Table.Head>Version</Table.Head>
						<Table.Head>Date</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					<Table.Row>
						<Table.Cell>{machine.hardwareInformation.mainboard_vendor}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.mainboard_product}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.mainboard_firmware_vendor}</Table.Cell>
						<Table.Cell>{machine.biosBootMethod.toUpperCase()}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.mainboard_firmware_version}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.mainboard_firmware_date}</Table.Cell>
					</Table.Row>
				</Table.Body>
			</Table.Root>
		</fieldset>
		<fieldset>
			<legend>Chassis</legend>
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>Vendor</Table.Head>
						<Table.Head>Type</Table.Head>
						<Table.Head>Version</Table.Head>
						<Table.Head>Serial</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					<Table.Row>
						<Table.Cell>{machine.hardwareInformation.chassis_vendor}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.chassis_type}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.chassis_version}</Table.Cell>
						<Table.Cell>{machine.hardwareInformation.chassis_serial}</Table.Cell>
					</Table.Row>
				</Table.Body>
			</Table.Root>
		</fieldset>
	</div>
{/snippet}
{#snippet TabContent_BlockDevices()}
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Name</Table.Head>
				<Table.Head>Model</Table.Head>
				<Table.Head>Serial</Table.Head>
				<Table.Head>Boot Disk</Table.Head>
				<Table.Head>Firmware Version</Table.Head>
				<Table.Head>Type</Table.Head>
				<Table.Head>User For</Table.Head>
				<Table.Head>Tags</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each machine.blockDevices as blockDevice}
				<Table.Row>
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
			<Table.Row>
				<Table.Head>
					<div class="text-sm">
						NAME
						<div class="text-xs font-light text-muted-foreground">MAC Address</div>
					</div>
				</Table.Head>
				<Table.Head>
					<div class="text-sm">
						IP ADDRESS
						<div class="text-xs font-light text-muted-foreground">Subnet</div>
					</div>
				</Table.Head>
				<Table.Head>
					<div class="text-sm">
						LINK SPEED
						<div class="text-xs font-light text-muted-foreground">Link Connected</div>
					</div>
				</Table.Head>
				<Table.Head>
					<div class="text-sm">
						FABRIC
						<div class="text-xs font-light text-muted-foreground">VLAN</div>
					</div>
				</Table.Head>
				<Table.Head>TYPE</Table.Head>
				<Table.Head>DHCP ON</Table.Head>
				<Table.Head>BOOT INTERFACE</Table.Head>
				<Table.Head>INTERFACE SPEED</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each machine.networkInterfaces as networkInterface}
				<Table.Row>
					<Table.Cell>
						{networkInterface.name}
						<div>
							{networkInterface.macAddress}
						</div>
					</Table.Cell>
					<Table.Cell>
						{networkInterface.ipAddress}
						<div>
							{networkInterface.subnetName}
						</div>
					</Table.Cell>
					<Table.Cell>
						{networkInterface.linkSpeed} Mbps
						<div>
							<Icon
								icon={networkInterface.linkConnected ? 'ph:check-circle' : 'ph:x-circle'}
								style="color: {networkInterface.linkConnected ? 'green' : 'red'}"
							/>
						</div>
					</Table.Cell>
					<Table.Cell>
						{networkInterface.fabricName}
						<div>
							{networkInterface.vlanName}
						</div>
					</Table.Cell>
					<Table.Cell>{networkInterface.type}</Table.Cell>
					<Table.Cell>
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
