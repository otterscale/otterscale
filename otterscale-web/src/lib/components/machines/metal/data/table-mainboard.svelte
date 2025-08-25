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
			<Table.Head>{m.product()}</Table.Head>
			<Table.Head>{m.firmware()}</Table.Head>
			<Table.Head>{m.boot_mode()}</Table.Head>
			<Table.Head>{m.version()}</Table.Head>
			<Table.Head>{m.date()}</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		<Table.Row>
			<Table.Cell>{$machine.hardwareInformation.mainboard_vendor}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.mainboard_product}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.mainboard_firmware_vendor}</Table.Cell>
			<Table.Cell>{$machine.biosBootMethod.toUpperCase()}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.mainboard_firmware_version}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.mainboard_firmware_date}</Table.Cell>
		</Table.Row>
	</Table.Body>
</Table.Root> -->

<Card.Root class="w-full max-w-sm shadow-none">
	<Card.Header class="bg-muted -mt-6 rounded-t-lg p-2">
		<Card.Title class="text-center">{m.mainboard()}</Card.Title>
	</Card.Header>
	<Card.Content>
		<div class="flex flex-col gap-6 text-sm">
			<div class="flex gap-2">
				<Icon icon="ph:factory" />
				<div>
					<Label>{m.vendor()}</Label>
					<p>{$machine.hardwareInformation.mainboard_vendor}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Icon icon="ph:barcode" />
				<div>
					<Label>{m.product()}</Label>
					<p>{$machine.hardwareInformation.mainboard_product}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Icon icon="ph:cpu" />
				<div>
					<Label>{m.firmware()}</Label>
					<p>{$machine.hardwareInformation.mainboard_firmware_vendor}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Icon icon="ph:power" />
				<div>
					<Label>{m.boot_mode()}</Label>
					<p>{$machine.biosBootMethod.toUpperCase()}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Icon icon="ph:tag" />
				<div>
					<Label>{m.version()}</Label>
					<p>{$machine.hardwareInformation.mainboard_firmware_version}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Icon icon="ph:calendar" />
				<div>
					<Label>{m.date()}</Label>
					<p>{$machine.hardwareInformation.mainboard_firmware_date}</p>
				</div>
			</div>
		</div>
	</Card.Content>
</Card.Root>
