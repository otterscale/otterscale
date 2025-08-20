<script lang="ts">
	import type { Network_IPAddress } from '$lib/api/network/v1/network_pb';
	import * as Table from '$lib/components/custom/table';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import Icon from '@iconify/svelte';

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
						IP
						<Table.SubHead>SYSTEM ID</Table.SubHead>
					</Table.Head>

					<Table.Head>
						HOST NAME
						<Table.SubHead>USER</Table.SubHead>
					</Table.Head>

					<Table.Head>
						TYPE
						<Table.SubHead>NODE TYPE</Table.SubHead>
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
							<Table.SubCell>
								{ipAddress.machineId}
							</Table.SubCell>
						</Table.Cell>

						<Table.Cell>
							{ipAddress.hostname}
							<Table.SubCell>
								{ipAddress.user}
							</Table.SubCell>
						</Table.Cell>

						<Table.Cell>
							{#if ipAddress.type}
								{ipAddress.type}
							{/if}
							<Table.SubCell>
								{#if ipAddress.nodeType}
									{ipAddress.nodeType}
								{/if}
							</Table.SubCell>
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</HoverCard.Content>
</HoverCard.Root>
