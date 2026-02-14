<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { Network_VLAN } from '$lib/api/network/v1/network_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { vlan }: { vlan: Network_VLAN } = $props();
</script>

<HoverCard.Root>
	<HoverCard.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
		<Icon icon="ph:info" />
	</HoverCard.Trigger>
	<HoverCard.Content class="w-fit">
		<Form.Root>
			<Form.Fieldset class="border-none p-2">
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General type="text" value={vlan.name} disabled />
				</Form.Field>

				{#if vlan.description}
					<Form.Field>
						<Form.Label>{m.description()}</Form.Label>
						<SingleInput.General type="text" value={vlan.description} disabled />
					</Form.Field>
				{/if}

				<Form.Field>
					<Form.Label>{m.mtu()}</Form.Label>
					<SingleInput.General type="number" value={vlan.mtu} disabled />
				</Form.Field>

				<Form.Field>
					<SingleInput.Boolean descriptor={() => m.dhcp_on()} value={vlan.dhcpOn} disabled />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
	</HoverCard.Content>
</HoverCard.Root>
