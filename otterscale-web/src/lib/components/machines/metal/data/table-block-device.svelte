<script lang="ts">
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Table from '$lib/components/custom/table';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { type Writable } from 'svelte/store';

	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<Table.Root>
	<Table.Header>
		<Table.Row>
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
		{#each $machine.blockDevices as blockDevice}
			<Table.Row>
				<Table.Cell>{blockDevice.name}</Table.Cell>
				<Table.Cell>{blockDevice.model}</Table.Cell>
				<Table.Cell>{blockDevice.serial}</Table.Cell>
				<Table.Cell>{blockDevice.bootDisk}</Table.Cell>
				<Table.Cell>{blockDevice.firmwareVersion}</Table.Cell>
				<Table.Cell>{blockDevice.type}</Table.Cell>
				<Table.Cell>{blockDevice.usedFor}</Table.Cell>
				<Table.Cell>
					{#each blockDevice.tags as tag}
						<Badge variant="outline">
							{tag}
						</Badge>
					{/each}
				</Table.Cell>
			</Table.Row>
		{/each}
	</Table.Body>
</Table.Root>
