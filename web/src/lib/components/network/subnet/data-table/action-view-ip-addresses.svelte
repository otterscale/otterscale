<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { Network_IPAddress } from '$lib/api/network/v1/network_pb';
	import { SubCell, SubHead } from '$lib/components/custom/table';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import * as Table from '$lib/components/ui/table';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { ipAddresses }: { ipAddresses: Network_IPAddress[] } = $props();
</script>

<HoverCard.Root>
	<HoverCard.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
		<Icon icon="ph:info" />
	</HoverCard.Trigger>
	<HoverCard.Content class="max-h-[50vh] w-fit overflow-y-auto">
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head></Table.Head>
					<Table.Head>
						{m.ip()}
						<SubHead>{m.system_id()}</SubHead>
					</Table.Head>

					<Table.Head>
						{m.hostname()}
						<SubHead>{m.user()}</SubHead>
					</Table.Head>

					<Table.Head>
						{m.type()}
						<SubHead>{m.node_type()}</SubHead>
					</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each ipAddresses as ipAddress, index}
					<Table.Row>
						<Table.Cell>
							{index + 1}
						</Table.Cell>

						<Table.Cell>
							{ipAddress.ip}
							<SubCell>
								{ipAddress.machineId}
							</SubCell>
						</Table.Cell>

						<Table.Cell>
							{ipAddress.hostname}
							<SubCell>
								{ipAddress.user}
							</SubCell>
						</Table.Cell>

						<Table.Cell>
							{#if ipAddress.type}
								{ipAddress.type}
							{/if}
							<SubCell>
								{#if ipAddress.nodeType}
									{ipAddress.nodeType}
								{/if}
							</SubCell>
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</HoverCard.Content>
</HoverCard.Root>
