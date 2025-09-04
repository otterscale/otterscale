<script lang="ts" module>
	import { type Writable } from 'svelte/store';

	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Table from '$lib/components/custom/table';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		machine,
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<Table.Root>
	<Table.Header>
		<Table.Row>
			<Table.Head>{m.name()}</Table.Head>
			<Table.Head>{m.machine_model()}</Table.Head>
			<Table.Head>{m.serial()}</Table.Head>
			<Table.Head>{m.boot_disk()}</Table.Head>
			<Table.Head>{m.firmware_version()}</Table.Head>
			<Table.Head>{m.type()}</Table.Head>
			<Table.Head>{m.used_for()}</Table.Head>
			<Table.Head>{m.tags()}</Table.Head>
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
