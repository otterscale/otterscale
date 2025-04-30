<script lang="ts">
	import CreateReservedIPRange from './create.svelte';
	import UpdateReservedIPRange from './update.svelte';
	import DeleteReservedIPRange from './delete.svelte';

	import Icon from '@iconify/svelte';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table/index.js';

	import { type Network_Subnet } from '$gen/api/nexus/v1/nexus_pb';

	let {
		subnet
	}: {
		subnet: Network_Subnet;
	} = $props();
</script>

<Sheet.Root>
	<Sheet.Trigger>
		<div class="flex items-center gap-1">
			<Icon icon="ph:arrow-square-out" class="size-4" />
		</div>
	</Sheet.Trigger>
	<Sheet.Content class="w-fit min-w-[50vw] overflow-auto" side="right">
		<Sheet.Header>
			<Sheet.Title>Reserved IP Ranges</Sheet.Title>
			<div class="ml-auto">
				<CreateReservedIPRange {subnet} />
			</div>
		</Sheet.Header>
		<Table.Root>
			<Table.Header class="sticky top-0 z-10 bg-background">
				<Table.Row>
					<Table.Head class="whitespace-nowrap font-light">ID</Table.Head>
					<Table.Head class="whitespace-nowrap font-light">TYPE</Table.Head>
					<Table.Head class="whitespace-nowrap font-light">START IP</Table.Head>
					<Table.Head class="whitespace-nowrap font-light">END IP</Table.Head>
					<Table.Head class="whitespace-nowrap font-light">COMMENT</Table.Head>
					<Table.Head class="whitespace-nowrap font-light"></Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body class="text-xs">
				{#each subnet.ipRanges.sort((previous, present) => Number(previous.id) - Number(present.id)) as ipRange}
					<Table.Row>
						<Table.Cell class="whitespace-nowrap">{ipRange.id}</Table.Cell>
						<Table.Cell class="whitespace-nowrap">
							{#if ipRange.type}
								<Badge variant="outline">{ipRange.type}</Badge>
							{/if}
						</Table.Cell>
						<Table.Cell class="whitespace-nowrap">{ipRange.startIp}</Table.Cell>
						<Table.Cell class="whitespace-nowrap">{ipRange.endIp}</Table.Cell>
						<Table.Cell class="whitespace-nowrap">{ipRange.comment}</Table.Cell>
						<Table.Cell class="flex justify-end whitespace-nowrap">
							<DropdownMenu.Root>
								<DropdownMenu.Trigger>
									<Button variant="ghost">
										<Icon icon="ph:dots-three-vertical" />
									</Button>
								</DropdownMenu.Trigger>
								<DropdownMenu.Content>
									<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
										<UpdateReservedIPRange {ipRange} />
									</DropdownMenu.Item>
									<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
										<DeleteReservedIPRange {ipRange} />
									</DropdownMenu.Item>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</Sheet.Content>
</Sheet.Root>
