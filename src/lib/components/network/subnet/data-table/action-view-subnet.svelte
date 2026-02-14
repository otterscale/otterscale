<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { Network_Subnet } from '$lib/api/network/v1/network_pb';
	import * as Form from '$lib/components/custom/form';
	import { Multiple as MultipleInput, Single as SingleInput } from '$lib/components/custom/input';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { subnet }: { subnet: Network_Subnet } = $props();
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
					<SingleInput.General type="text" value={subnet.name} disabled />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.cidr()}</Form.Label>
					<SingleInput.General type="text" value={subnet.cidr} disabled />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.gateway()}</Form.Label>
					<SingleInput.General type="text" value={subnet.gatewayIp} disabled />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.dns_server()}</Form.Label>
					<MultipleInput.Root type="text" values={subnet.dnsServers}>
						<MultipleInput.Viewer disabled />
					</MultipleInput.Root>
				</Form.Field>

				{#if subnet.description}
					<Form.Field>
						<Form.Label>{m.description()}</Form.Label>
						<SingleInput.General type="text" value={subnet.description} disabled />
					</Form.Field>
				{/if}

				<Form.Field>
					<SingleInput.Boolean
						descriptor={() => m.managed_allocation()}
						value={subnet.managedAllocation}
						disabled
					/>
				</Form.Field>

				<Form.Field>
					<SingleInput.Boolean
						descriptor={() => m.active_discovery()}
						value={subnet.activeDiscovery}
						disabled
					/>
				</Form.Field>

				<Form.Field>
					<SingleInput.Boolean
						descriptor={() => m.allow_dns_resolution()}
						value={subnet.allowDnsResolution}
						disabled
					/>
				</Form.Field>

				<Form.Field>
					<SingleInput.Boolean
						descriptor={() => m.allow_proxy_access()}
						value={subnet.allowProxyAccess}
						disabled
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
	</HoverCard.Content>
</HoverCard.Root>
