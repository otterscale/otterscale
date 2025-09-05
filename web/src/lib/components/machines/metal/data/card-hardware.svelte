<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';

	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import Label from '$lib/components/ui/label/label.svelte';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	import Separator from '$lib/components/ui/separator/separator.svelte';

	let {
		machine,
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<div class="flex gap-8 p-4">
	<div class="flex flex-col gap-4">
		<div class="flex items-center gap-2 text-center text-lg">
			<Icon icon="ph:desktop" class="size-5" />
			{m.system()}
		</div>
		<Separator />
		<div class="flex flex-col gap-4 text-sm">
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:cpu" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.hardwareInformation.system_vendor}</p>
					<Label class="text-muted-foreground">{m.vendor()}</Label>
				</div>
			</div>
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:archive-box" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.hardwareInformation.system_product}</p>
					<Label class="text-muted-foreground">{m.product()}</Label>
				</div>
			</div>
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:tag" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.hardwareInformation.system_version}</p>
					<Label class="text-muted-foreground">{m.version()}</Label>
				</div>
			</div>
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:barcode" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.hardwareInformation.system_serial}</p>
					<Label class="text-muted-foreground">{m.serial()}</Label>
				</div>
			</div>
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:ticket" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.hardwareInformation.system_sku}</p>
					<Label class="text-muted-foreground">{m.sku()}</Label>
				</div>
			</div>
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:users-three" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.hardwareInformation.system_family}</p>
					<Label class="text-muted-foreground">{m.family()}</Label>
				</div>
			</div>
		</div>
	</div>
	<div class="flex flex-col gap-4">
		<div class="flex items-center gap-2 text-center text-lg">
			<Icon icon="ph:circuitry" class="size-5" />
			{m.mainboard()}
		</div>
		<Separator />
		<div class="flex flex-col gap-4 text-sm">
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:factory" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.hardwareInformation.mainboard_vendor}</p>
					<Label class="text-muted-foreground">{m.vendor()}</Label>
				</div>
			</div>
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:barcode" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.hardwareInformation.mainboard_product}</p>
					<Label class="text-muted-foreground">{m.product()}</Label>
				</div>
			</div>
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:cpu" class="size-4" />
				<div>
					<p class="whitespace-nowrap">
						{$machine.hardwareInformation.mainboard_firmware_vendor}
					</p>
					<Label class="text-muted-foreground">{m.firmware()}</Label>
				</div>
			</div>
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:power" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.biosBootMethod.toUpperCase()}</p>
					<Label class="text-muted-foreground">{m.boot_mode()}</Label>
				</div>
			</div>
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:tag" class="size-4" />
				<div>
					<p class="whitespace-nowrap">
						{$machine.hardwareInformation.mainboard_firmware_version}
					</p>
					<Label class="text-muted-foreground">{m.version()}</Label>
				</div>
			</div>
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:calendar" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.hardwareInformation.mainboard_firmware_date}</p>
					<Label class="text-muted-foreground">{m.date()}</Label>
				</div>
			</div>
		</div>
	</div>
	<div class="flex flex-col gap-4">
		<div class="flex items-center gap-2 text-center text-lg">
			<Icon icon="ph:cube" class="size-5" />
			{m.chassis()}
		</div>
		<Separator />
		<div class="flex flex-col gap-4 text-sm">
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:factory" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.hardwareInformation.chassis_vendor}</p>
					<Label class="text-muted-foreground">{m.vendor()}</Label>
				</div>
			</div>
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:archive-box" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.hardwareInformation.chassis_type}</p>
					<Label class="text-muted-foreground">{m.type()}</Label>
				</div>
			</div>
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:tag" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.hardwareInformation.chassis_version}</p>
					<Label class="text-muted-foreground">{m.version()}</Label>
				</div>
			</div>
			<div class="flex items-center gap-2 pl-1">
				<Icon icon="ph:barcode" class="size-4" />
				<div>
					<p class="whitespace-nowrap">{$machine.hardwareInformation.chassis_serial}</p>
					<Label class="text-muted-foreground">{m.serial()}</Label>
				</div>
			</div>
		</div>
	</div>
</div>
