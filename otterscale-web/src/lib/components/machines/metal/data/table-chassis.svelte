<script lang="ts" module>
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Card from '$lib/components/ui/card';
	import Label from '$lib/components/ui/label/label.svelte';
	import { m } from '$lib/paraglide/messages';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<!-- 
<Table.Root>
	<Table.Header>
		<Table.Row>
			<Table.Head>{m.vendor()}</Table.Head>
			<Table.Head>{m.type()}</Table.Head>
			<Table.Head>{m.version()}</Table.Head>
			<Table.Head>{m.serial()}</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		<Table.Row>
			<Table.Cell>{$machine.hardwareInformation.chassis_vendor}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.chassis_type}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.chassis_version}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.chassis_serial}</Table.Cell>
		</Table.Row>
	</Table.Body>
</Table.Root> -->

<Card.Root class="w-full max-w-sm shadow-none">
	<Card.Header class="bg-muted -mt-6 rounded-t-lg p-2">
		<Card.Title class="text-center">{m.chassis()}</Card.Title>
	</Card.Header>
	<Card.Content>
		<div class="flex flex-col gap-6 text-sm">
			<div class="flex gap-2">
				<Icon icon="ph:factory" />
				<div>
					<Label>{m.vendor()}</Label>
					<p>{$machine.hardwareInformation.chassis_vendor}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Icon icon="ph:archive-box" />
				<div>
					<Label>{m.type()}</Label>
					<p>{$machine.hardwareInformation.chassis_type}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Icon icon="ph:tag" />
				<div>
					<Label>{m.version()}</Label>
					<p>{$machine.hardwareInformation.chassis_version}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Icon icon="ph:barcode" />
				<div>
					<Label>{m.serial()}</Label>
					<p>{$machine.hardwareInformation.chassis_serial}</p>
				</div>
			</div>
		</div>
	</Card.Content>
</Card.Root>
