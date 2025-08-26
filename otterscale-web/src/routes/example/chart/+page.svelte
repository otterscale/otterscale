<script lang="ts">
	import Layout from '$lib/components/custom/chart/layout/small-flexible-height.svelte';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { formatBigNumber as formatNumber, formatProgressColor } from '$lib/formatter';
	import Title from '$lib/components/custom/chart/title.svelte';
	import * as Card from '$lib/components/ui/card/index.js';

	const cpu = 22;
	const ram = 70;
	const disk = 40;
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	<Card.Root>
		<Card.Header class="h-[30px] items-center">
			<Card.Title>USAGE RATE</Card.Title>
		</Card.Header>

		<Card.Content class="flex h-[250px] flex-col space-y-4">
			<div class="space-y-2">
				<div class="flex items-center justify-between">
					<span class="text-muted-foreground text-sm">CPU</span>
					<span class="text-muted-foreground text-sm">{cpu}%</span>
				</div>
				<Progress value={cpu} max={100} class={formatProgressColor(cpu, true, false)} />
			</div>

			<div class="space-y-2">
				<div class="flex items-center justify-between">
					<span class="text-muted-foreground text-sm">RAM</span>
					<span class="text-muted-foreground text-sm">{ram}%</span>
				</div>
				<Progress value={ram} max={100} class={formatProgressColor(ram, true, false)} />
			</div>

			<div class="space-y-2">
				<div class="flex items-center justify-between">
					<span class="text-muted-foreground text-sm">DISK</span>
					<span class="text-muted-foreground text-sm">{disk}%</span>
				</div>
				<Progress value={disk} max={100} class={formatProgressColor(disk, true, false)} />
			</div>
		</Card.Content>
	</Card.Root>

	<Layout>
		{#snippet title()}
			<Title title="AVAILABILITY" />
		{/snippet}

		{#snippet content()}
			<Progress value={22} max={100} class={formatProgressColor(22)} />
		{/snippet}
	</Layout>
</div>
