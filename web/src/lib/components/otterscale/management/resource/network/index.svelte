<script lang="ts">
	import { page } from '$app/state';
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
	import { createClient, type Transport } from '@connectrpc/connect';
	import { writable } from 'svelte/store';
	import { getContext, onMount } from 'svelte';
	import { formatBigNumber as formatNumber } from '$lib/formatter';

	import { ManagementNetworkSubnetReservedIPRanges } from '$lib/components/otterscale/index';

	import {
		Nexus,
		type Network,
		type Network_Subnet,
		type Network_VLAN,
		type Network_IPAddress,
		type Scope,
		type Tag
	} from '$gen/api/nexus/v1/nexus_pb';

	let {
		networks
	}: {
		networks: Network[];
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const networksStore = writable<Network[]>([]);
	const networksLoading = writable(true);
	async function refreshNetworks() {
		while (page.url.searchParams.get('intervals')) {
			await new Promise((resolve) =>
				setTimeout(resolve, 1000 * Number(page.url.searchParams.get('intervals')))
			);
			console.log(`Refresh networks`);

			try {
				const response = await client.listNetworks({});
				networksStore.set(response.networks);
			} catch (error) {
				console.error('Error fetching:', error);
			}
		}
	}

	let mounted = false;
	onMount(async () => {
		try {
			refreshNetworks();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<div>
	{@render StatisticNetworks()}
	<div class="p-4">
		<div class="flex justify-end py-2">
			<CreateNetwork bind:networks />
		</div>
		<Table.Root>
			<Table.Header class="bg-muted/50">
				<Table.Row class="*:text-xs *:font-light [&>th]:py-2">
					<Table.Head>FABRIC</Table.Head>
					<Table.Head>VLAN</Table.Head>
					<Table.Head>DHCP</Table.Head>
					<Table.Head>SUBNET</Table.Head>
					<Table.Head>USED IP</Table.Head>
					<Table.Head>RESERVED IP RANGE</Table.Head>
					<Table.Head>AVAILABLE</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each networks.sort((previous, prevent) => (previous.fabric?.name ?? '').localeCompare(prevent.fabric?.name ?? '')) as network}
					<Table.Row>
						{#if network.fabric}
							<Table.Cell>
								<div class="flex items-center justify-between">
									{network.fabric.name}
									<DropdownMenu.Root>
										<DropdownMenu.Trigger>
											<Button variant="ghost">
												<Icon icon="ph:dots-three-vertical" />
											</Button>
										</DropdownMenu.Trigger>
										<DropdownMenu.Content>
											<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
												<UpdateFabric bind:networks fabric={network.fabric} />
											</DropdownMenu.Item>
											<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
												<DeleteNetwork bind:networks fabric={network.fabric} />
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
													<UpdateVLAN bind:networks fabric={network.fabric} vlan={network.vlan} />
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
													<UpdateSubnet bind:networks subnet={network.subnet} />
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

	<div class="grid grid-cols-4 gap-3">
		<Card.Root>
			<Card.Header class="h-10">
				<Card.Title>NETWORK</Card.Title>
			</Card.Header>
			<Card.Content class="h-30">
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
			<Card.Header class="h-10">
				<Card.Title>DHCP</Card.Title>
			</Card.Header>
			<Card.Content class="h-30">
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
			<Card.Header class="h-10">
				<Card.Title>AVAILABILITY</Card.Title>
			</Card.Header>
			<Card.Content class="h-30">
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
			<Table.Root class="*:whitespace-nowrap *:text-xs">
				<Table.Body class="[&>tr>th]:text-right [&>tr>th]:font-light">
					<Table.Row>
						<Table.Head>Name</Table.Head>
						<Table.Cell>{vlan.name}</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>MTU</Table.Head>
						<Table.Cell>{vlan.mtu}</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Description</Table.Head>
						<Table.Cell>{vlan.description}</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>DHCP</Table.Head>
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
		<HoverCard.Content class="w-fit">
			<Table.Root class="*:whitespace-nowrap *:text-xs">
				<Table.Body class="[&>tr>th]:text-right [&>tr>th]:font-light">
					<Table.Row>
						<Table.Head>Name</Table.Head>
						<Table.Cell>
							{subnet.name}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>CIDR</Table.Head>
						<Table.Cell>
							{subnet.cidr}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Gateway IP</Table.Head>
						<Table.Cell>
							{subnet.gatewayIp}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>DNS Servers</Table.Head>
						<Table.Cell>
							{#each subnet.dnsServers as dnsServer}
								<Badge variant="outline">
									{dnsServer}
								</Badge>
							{/each}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Description</Table.Head>
						<Table.Cell>
							{subnet.description}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Managed Allocation</Table.Head>
						<Table.Cell>
							{@render formatterBoolean(subnet.managedAllocation)}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Active Discovery</Table.Head>
						<Table.Cell>
							{@render formatterBoolean(subnet.activeDiscovery)}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Allow Proxy Access</Table.Head>
						<Table.Cell>
							{@render formatterBoolean(subnet.allowProxyAccess)}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Allow DNS Resolution</Table.Head>
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
			<HoverCard.Content class="max-h-[50vh] min-w-[50vw] overflow-auto">
				<Table.Root class="w-full *:whitespace-nowrap *:text-xs">
					<Table.Header>
						<Table.Row class="*:font-light">
							<Table.Head></Table.Head>
							<Table.Head>IP</Table.Head>
							<Table.Head>User</Table.Head>
							<Table.Head>Node Type</Table.Head>
							<Table.Head>Type</Table.Head>
							<Table.Head>System ID</Table.Head>
							<Table.Head>Hostname</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each ipAddresses as ipAddress, index}
							<Table.Row>
								<Table.Cell>{index + 1}</Table.Cell>
								<Table.Cell>{ipAddress.ip}</Table.Cell>
								<Table.Cell>
									{ipAddress.user}
								</Table.Cell>
								<Table.Cell>
									<span>
										{#if ipAddress.nodeType}
											<Badge variant="outline">{ipAddress.nodeType}</Badge>
										{/if}
									</span>
								</Table.Cell>
								<Table.Cell>
									<span>
										{#if ipAddress.type}
											<Badge variant="outline">{ipAddress.type}</Badge>
										{/if}
									</span>
								</Table.Cell>
								<Table.Cell>{ipAddress.ip}</Table.Cell>
								<Table.Cell>{ipAddress.hostname}</Table.Cell>
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
