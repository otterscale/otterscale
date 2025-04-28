<script lang="ts">
	import CreateNetwork from './create.svelte';
	import DeleteNetwork from './delete.svelte';
	import UpdateFabric from './update-fabric.svelte';
	import UpdateSubnet from './update-subnet.svelte';
	import UpdateVLAN from './update-vlan.svelte';

	import Icon from '@iconify/svelte';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';

	import { writable } from 'svelte/store';

	import { formatBigNumber as formatNumber } from '$lib/formatter';

	import { ManagementNetworkSubnetReservedIPRanges } from '$lib/components/otterscale/index';

	import {
		type Network,
		type Network_Subnet,
		type Network_VLAN,
		type Network_IPAddress,
		type Scope,
		type Tag
	} from '$gen/api/nexus/v1/nexus_pb';
	import { onMount } from 'svelte';

	let {
		networks
	}: {
		networks: Network[];
	} = $props();
</script>

<div class="space-y-3">
	{@render StatisticNetworks()}
	<span class="flex justify-end">
		<CreateNetwork />
	</span>
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head class="text-xs font-light">FABRIC</Table.Head>
				<Table.Head class="text-xs font-light">VLAN</Table.Head>
				<Table.Head class="text-xs font-light">DHCP</Table.Head>
				<Table.Head class="text-xs font-light">SUBNET</Table.Head>
				<Table.Head class="text-xs font-light">USED IP</Table.Head>
				<Table.Head class="text-xs font-light">RESERVED IP RANGE</Table.Head>
				<Table.Head class="text-xs font-light">AVAILABLE</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each networks as network}
				<Table.Row>
					{#if network.fabric}
						<Table.Cell>
							<div class="flex justify-between">
								{network.fabric.name}
								<DropdownMenu.Root>
									<DropdownMenu.Trigger>
										<Button variant="ghost">
											<Icon icon="ph:dots-three-vertical" />
										</Button>
									</DropdownMenu.Trigger>
									<DropdownMenu.Content>
										<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
											<UpdateFabric fabric={network.fabric} />
										</DropdownMenu.Item>
										<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
											<DeleteNetwork fabric={network.fabric} />
										</DropdownMenu.Item>
									</DropdownMenu.Content>
								</DropdownMenu.Root>
							</div>
						</Table.Cell>
						<Table.Cell>
							{#if network.vlan}
								<span class="flex justify-between">
									<span class="flex items-center gap-1">
										{network.vlan.name}
										{@render ReadVLAN(network.vlan)}
									</span>
									<DropdownMenu.Root>
										<DropdownMenu.Trigger>
											<Button variant="ghost">
												<Icon icon="ph:dots-three-vertical" />
											</Button>
										</DropdownMenu.Trigger>
										<DropdownMenu.Content>
											<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
												<UpdateVLAN fabric={network.fabric} vlan={network.vlan} />
											</DropdownMenu.Item>
										</DropdownMenu.Content>
									</DropdownMenu.Root>
								</span>
							{/if}
						</Table.Cell>
						<Table.Cell>
							{#if network.vlan}
								{@render formatterBoolean(network.vlan.dhcpOn)}
							{/if}
						</Table.Cell>
						<Table.Cell>
							{#if network.subnet}
								<div class="flex justify-between">
									<span class="flex items-center gap-1">
										{network.subnet.name}
										{@render ReadSubnet(network.subnet)}
									</span>
									<DropdownMenu.Root>
										<DropdownMenu.Trigger>
											<Button variant="ghost">
												<Icon icon="ph:dots-three-vertical" />
											</Button>
										</DropdownMenu.Trigger>
										<DropdownMenu.Content>
											<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
												<UpdateSubnet subnet={network.subnet} />
											</DropdownMenu.Item>
										</DropdownMenu.Content>
									</DropdownMenu.Root>
								</div>
							{/if}
						</Table.Cell>
						<Table.Cell>
							{#if network.subnet}
								<span class="flex items-center gap-1">
									{network.subnet.ipAddresses.length}
									{@render ReadUsedIP(network.subnet.ipAddresses)}
								</span>
							{/if}
						</Table.Cell>
						<Table.Cell>
							{#if network.subnet}
								<span class="flex items-center gap-1">
									{network.subnet.ipRanges.length}
									<ManagementNetworkSubnetReservedIPRanges subnet={network.subnet} />
								</span>
							{/if}
						</Table.Cell>
						<Table.Cell>
							{#if network.subnet && network.subnet.statistics}
								{@const availabilityBySubnet =
									(Number(network.subnet.statistics.available) * 100) /
									Number(network.subnet.statistics.total)}

								<div class="flex justify-between">
									<p>{network.subnet.statistics.availablePercent}</p>
									<p class="text-xs font-extralight">
										{formatNumber(network.subnet.statistics.available)}/{formatNumber(
											network.subnet.statistics.total
										)}
									</p>
								</div>
								<Progress
									max={100}
									value={availabilityBySubnet}
									class={`${
										availabilityBySubnet > 62
											? ' bg-green-50 *:bg-green-700'
											: availabilityBySubnet > 38
												? 'bg-yellow-50 *:bg-yellow-500'
												: 'bg-red-50 *:bg-red-700'
									}`}
								/>
							{/if}
						</Table.Cell>
					{/if}
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>

{#snippet StatisticNetworks()}
	{@const numberOfFabrics = networks.filter((n) => n).map((n) => n.fabric).length}
	{@const vlans = networks.filter((n) => n).map((n) => n.vlan)}
	{@const numberOfDHCPOn = networks.reduce((a, network) => a + (network.vlan?.dhcpOn ? 1 : 0), 0)}
	{@const numberOfAvailableSubnet = networks.reduce(
		(a, network) =>
			a + Number(network.subnet?.statistics ? network.subnet.statistics.available : 0),
		0
	)}
	{@const numberOfSubnet = networks.reduce(
		(a, network) => a + Number(network.subnet?.statistics ? network.subnet.statistics.total : 0),
		0
	)}
	{@const Availability = (numberOfAvailableSubnet * 100) / numberOfSubnet || 0}
	{@const rateOfDHCPOn = (numberOfDHCPOn * 100) / numberOfFabrics || 0}
	<div class="grid grid-cols-4 gap-3 *:border-none *:shadow-none">
		<Card.Root>
			<Card.Header>
				<Card.Title>Networks</Card.Title>
			</Card.Header>
			<Card.Content>
				<span class="flex items-end gap-1">
					<p class="text-7xl">{numberOfFabrics}</p>
					<p>fabrics</p>
				</span>
			</Card.Content>
			<Card.Footer>
				<div class="flex flex-wrap gap-1">
					{#each [...new Set(vlans.map((v) => v?.name))] as vlanName}
						<Badge variant="outline">
							{vlanName}: {networks.reduce(
								(a, network) => a + (network.vlan?.name === vlanName ? 1 : 0),
								0
							)}
						</Badge>
					{/each}
				</div>
			</Card.Footer>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title>DHCP</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-3xl">
					{Math.round(rateOfDHCPOn)}%
				</p>
				<p class="text-xs text-muted-foreground">
					{numberOfDHCPOn} Enabled over {numberOfFabrics} fabrics
				</p>
			</Card.Content>
			<Card.Footer>
				<Progress value={rateOfDHCPOn} max={100} />
			</Card.Footer>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title>Availability</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-3xl">
					{Math.round(Availability)}%
				</p>
				<p class="text-xs text-muted-foreground">
					{formatNumber(numberOfAvailableSubnet)} available over {formatNumber(numberOfSubnet)} units
				</p>
			</Card.Content>
			<Card.Footer>
				<Progress
					value={Availability}
					max={100}
					class={`${
						Availability > 62
							? 'bg-green-50 *:bg-green-700'
							: Availability > 38
								? 'bg-yellow-50 *:bg-yellow-500'
								: 'bg-red-50 *:bg-red-700'
					}`}
				/>
			</Card.Footer>
		</Card.Root>
	</div>
{/snippet}

{#snippet ReadVLAN(vlan: Network_VLAN)}
	<HoverCard.Root openDelay={13}>
		<HoverCard.Trigger>
			<div class="flex items-center gap-1">
				<Icon icon="ph:info" class="size-4 text-blue-800" />
			</div>
		</HoverCard.Trigger>
		<HoverCard.Content>
			<Table.Root>
				<Table.Body class="text-xs">
					<Table.Row>
						<Table.Cell>Name</Table.Cell>
						<Table.Cell>{vlan.name}</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell>MTU</Table.Cell>
						<Table.Cell>{vlan.mtu}</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell>Description</Table.Cell>
						<Table.Cell>{vlan.description}</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell>DHCP</Table.Cell>
						<Table.Cell>
							{@render formatterBoolean(vlan.dhcpOn)}
						</Table.Cell>
					</Table.Row>
				</Table.Body>
			</Table.Root>
		</HoverCard.Content>
	</HoverCard.Root>
{/snippet}

{#snippet ReadSubnet(subnet: Network_Subnet)}
	<HoverCard.Root openDelay={13}>
		<HoverCard.Trigger>
			<div class="flex items-center gap-1">
				<Icon icon="ph:info" class="size-4 text-blue-800" />
			</div>
		</HoverCard.Trigger>
		<HoverCard.Content>
			<Table.Root>
				<Table.Body class="text-xs">
					<Table.Row>
						<Table.Cell>Name</Table.Cell>
						<Table.Cell>
							{subnet.name}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell>CIDR</Table.Cell>
						<Table.Cell>
							{subnet.cidr}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell>Gateway IP</Table.Cell>
						<Table.Cell>
							{subnet.gatewayIp}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell>DNS Servers</Table.Cell>
						<Table.Cell>
							{#each subnet.dnsServers as dnsServer}
								<Badge class="text-xs" variant="outline">
									{dnsServer}
								</Badge>
							{/each}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell>Description</Table.Cell>
						<Table.Cell>
							{subnet.description}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell>Managed Allocation</Table.Cell>
						<Table.Cell>
							{@render formatterBoolean(subnet.managedAllocation)}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell>Active Discovery</Table.Cell>
						<Table.Cell>
							{@render formatterBoolean(subnet.activeDiscovery)}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell>Allow Proxy Access</Table.Cell>
						<Table.Cell>
							{@render formatterBoolean(subnet.allowProxyAccess)}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell>Allow DNS Resolution</Table.Cell>
						<Table.Cell>
							{@render formatterBoolean(subnet.allowDnsResolution)}
						</Table.Cell>
					</Table.Row>
				</Table.Body>
			</Table.Root>
		</HoverCard.Content>
	</HoverCard.Root>
{/snippet}

{#snippet ReadUsedIP(ipAddresses: Network_IPAddress[])}
	{#if ipAddresses.length > 0}
		<HoverCard.Root openDelay={13}>
			<HoverCard.Trigger>
				<div class="flex items-center gap-1">
					<Icon icon="ph:info" class="size-4 text-blue-800" />
				</div>
			</HoverCard.Trigger>
			<HoverCard.Content class="max-h-[50vh] min-w-[38vw] max-w-[62vw] overflow-auto">
				<Table.Root class="w-fit">
					<Table.Header>
						<Table.Row>
							<Table.Head class="whitespace-nowrap font-light"></Table.Head>
							<Table.Head class="whitespace-nowrap font-light">IP</Table.Head>
							<Table.Head class="whitespace-nowrap font-light">USER</Table.Head>
							<Table.Head class="whitespace-nowrap font-light">NODE TYPE</Table.Head>
							<Table.Head class="whitespace-nowrap font-light">TYPE</Table.Head>
							<Table.Head class="whitespace-nowrap font-light">SYSTEM ID</Table.Head>
							<Table.Head class="whitespace-nowrap font-light">HOSTNAME</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body class="text-xs">
						{#each ipAddresses as ipAddress, index}
							<Table.Row>
								<Table.Cell class="whitespace-nowrap">{index + 1}</Table.Cell>
								<Table.Cell class="whitespace-nowrap">{ipAddress.ip}</Table.Cell>
								<Table.Cell class="whitespace-nowrap">
									{ipAddress.user}
								</Table.Cell>
								<Table.Cell class="whitespace-nowrap">
									{#if ipAddress.nodeType}
										<Badge variant="outline">{ipAddress.nodeType}</Badge>
									{/if}
								</Table.Cell>
								<Table.Cell class="whitespace-nowrap">
									{#if ipAddress.type}
										<Badge variant="outline">{ipAddress.type}</Badge>
									{/if}
								</Table.Cell>
								<Table.Cell class="whitespace-nowrap">{ipAddress.ip}</Table.Cell>
								<Table.Cell class="whitespace-nowrap">{ipAddress.hostname}</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</HoverCard.Content>
		</HoverCard.Root>
	{/if}
{/snippet}

{#snippet formatterBoolean(b: boolean)}
	<Icon icon={b ? 'ph:check' : 'ph:x'} />
{/snippet}
