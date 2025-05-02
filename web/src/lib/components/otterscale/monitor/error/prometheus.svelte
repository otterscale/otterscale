<script lang="ts">
	import Icon from '@iconify/svelte';
	import { cn } from '$lib/utils';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { Svg, Group, Arc, Chart, Text } from 'layerchart';
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';

	const chartHeight = {
		large: 150,
		small: 100
	};
</script>

<div class="grid w-full gap-2">
	<Alert.Root
		class="absolute z-10 flex w-[calc(100vw_-_theme(spacing.72))] justify-between bg-yellow-50 opacity-95"
	>
		<span class=" flex items-center gap-2">
			<Icon icon="radix-icons:exclamation-triangle" class="size-11 text-yellow-400" />
			<span>
				<Alert.Title class="font-bold text-yellow-500">WARNING</Alert.Title>
				<Alert.Description class="font-base text-yellow-500"
					>Prometheus is required for monitoring and metrics collection. Install it to enable system
					monitoring capabilities.</Alert.Description
				>
			</span>
		</span>
		<Button variant="destructive" class="text-sm">Install</Button>
	</Alert.Root>

	<div>
		<div class="col-span-1 *:animate-pulse">
			<div class="col-span-1 grid grid-cols-2">
				{#each Array(2) as _, i}
					<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
						<Card.Header>
							<Card.Title class="flex justify-between">
								<h1 class="text-3xl">Health</h1>
							</Card.Title>
							<Card.Description>
								<span class="text-sm font-light">Ratio of non-running pods over all pods</span>
							</Card.Description>
						</Card.Header>
						<Card.Content>
							<div class="flex h-full w-full items-center justify-center">
								<div
									class={cn(
										`h-[${Math.round(chartHeight.large * 1.732)}px] w-[${Math.round(chartHeight.large * 1.732)}px]`
									)}
								>
									<Chart>
										<Svg center>
											<Group y={chartHeight.large / 4}>
												<Arc
													value={0}
													domain={[0, 100]}
													outerRadius={chartHeight.large}
													innerRadius={-13}
													cornerRadius={13}
													range={[-120, 120]}
													track={{ class: 'fill-muted' }}
												>
													<Text
														value={`?`}
														textAnchor="middle"
														verticalAnchor="middle"
														class="text-5xl tabular-nums"
													/>
												</Arc>
											</Group>
										</Svg>
									</Chart>
								</div>
							</div>
						</Card.Content>
						<Card.Footer></Card.Footer>
					</Card.Root>
				{/each}
			</div>

			<div class="col-span-1 grid grid-cols-1 gap-2">
				<div class="col-span-1 grid grid-cols-3 gap-3">
					{#each Array(3) as _, i}
						<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
							<Card.Header>
								<Card.Title class="text-2xl">CPU</Card.Title>
								<Card.Description
									><span class="text-sm font-light">CPU usage across all clusters</span
									></Card.Description
								>
							</Card.Header>
							<Card.Content>
								<div class="flex h-full w-full items-center justify-center">
									<div
										class={cn(
											`h-[${Math.round(chartHeight.small * 1.732)}px] w-[${Math.round(chartHeight.small * 1.732)}px]`
										)}
									>
										<Chart>
											<Svg center>
												<Group y={chartHeight.small / 4}>
													<Arc
														value={0}
														domain={[0, 100]}
														outerRadius={chartHeight.small}
														innerRadius={-13}
														cornerRadius={13}
														range={[-120, 120]}
														track={{ class: 'fill-muted' }}
													>
														<Text
															value={`?`}
															textAnchor="middle"
															verticalAnchor="middle"
															class="text-5xl tabular-nums"
														/>
													</Arc>
												</Group>
											</Svg>
										</Chart>
									</div>
								</div>
							</Card.Content>
						</Card.Root>
					{/each}
				</div>
			</div>
		</div>
	</div>
</div>
