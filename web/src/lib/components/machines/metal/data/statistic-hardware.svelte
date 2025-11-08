<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';

	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Layout } from '$lib/components/custom/instance';
	import { buttonVariants } from '$lib/components/ui/button/button.svelte';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import { m } from '$lib/paraglide/messages';

	import CardHardware from './card-hardware.svelte';
</script>

<script lang="ts">
	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<Layout.Statistic.Root>
	<Layout.Statistic.Header>
		<Layout.Statistic.Title>{m.hardware_information()}</Layout.Statistic.Title>
		<Layout.Statistic.Action>
			<HoverCard.Root>
				<HoverCard.Trigger class={buttonVariants({ variant: 'ghost' })}>
					<Icon icon="ph:info" />
				</HoverCard.Trigger>
				<HoverCard.Content class="w-fit max-w-[77vw]">
					<CardHardware {machine} />
				</HoverCard.Content>
			</HoverCard.Root>
		</Layout.Statistic.Action>
	</Layout.Statistic.Header>
	<Layout.Statistic.Content>
		<div class="flex-col text-base">
			{$machine.hardwareInformation.system_product}
			<p class="text-muted-foreground">{$machine.hardwareInformation.system_vendor}</p>
		</div>
	</Layout.Statistic.Content>
	<Layout.Statistic.Footer>
		<div class="flex gap-2">
			<p>{$machine.hardwareInformation.mainboard_firmware_vendor}</p>
			<p>{$machine.hardwareInformation.mainboard_firmware_version}</p>
		</div>
	</Layout.Statistic.Footer>
</Layout.Statistic.Root>
