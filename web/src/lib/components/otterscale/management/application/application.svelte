<script lang="ts">
	import { type Application } from '$gen/api/nexus/v1/nexus_pb';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';

	import Icon from '@iconify/svelte';
	import { Badge } from '$lib/components/ui/badge';
	import * as Table from '$lib/components/ui/table';
	import * as Card from '$lib/components/ui/card';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';

	let {
		application
	}: {
		application: Application;
	} = $props();
</script>

{#snippet StatisticApplication(application: Application)}
	{@const numberOfPods = application.pods.length}
	{@const numberOfContainers = application.containers.length}
	{@const numberOfPersistentVolumeClaims = application.persistentVolumeClaims.length}
	{@const numberOfHealthyPods = application.healthies}
	{@const health = (numberOfHealthyPods * 100) / numberOfPods || 0}

	<Card.Root>
		<Card.Header>
			<Card.Title>Containers</Card.Title>
		</Card.Header>
		<Card.Content class="text-3xl">
			{numberOfContainers}
		</Card.Content>
		<Card.Footer class="flex flex-col items-start">
			{#each application.containers as container}
				<span class="flex items-center gap-1">
					<HoverCard.Root openDelay={13}>
						<HoverCard.Trigger>
							<Icon icon="ph:info" class="size-4 text-blue-800" />
						</HoverCard.Trigger>
						<HoverCard.Content class="min-w-fit">
							<Table.Root>
								<Table.Body>
									<Table.Row class="*:whitespace-nowrap">
										<Table.Head>Image</Table.Head>
										<Table.Cell>{container.imageName}</Table.Cell>
									</Table.Row>
									<Table.Row class="*:whitespace-nowrap">
										<Table.Head>Pull Policy</Table.Head>
										<Table.Cell>
											<Badge variant="outline">{container.imagePullPolicy}</Badge>
										</Table.Cell>
									</Table.Row>
								</Table.Body>
							</Table.Root>
						</HoverCard.Content>
					</HoverCard.Root>
					<Badge variant="outline">
						<!-- <p class=" overflow-visible whitespace-nowrap text-[13px]"> -->
						{container.imageName}
						<!-- </p> -->
					</Badge>
				</span>
			{/each}
		</Card.Footer>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Volumes</Card.Title>
		</Card.Header>
		<Card.Content class="text-3xl">
			{numberOfPersistentVolumeClaims}
		</Card.Content>
		<Card.Footer class="flex flex-col items-start">
			{#each application.persistentVolumeClaims as persistentVolumeClaim}
				<span class="flex items-center gap-1">
					<HoverCard.Root openDelay={13}>
						<HoverCard.Trigger>
							<Icon icon="ph:info" class="size-4 text-blue-800" />
						</HoverCard.Trigger>
						<HoverCard.Content class="min-w-fit">
							<fieldset>
								<legend class="text-center text-xs font-light">Persistent Volumn Claim</legend>
								<Table.Root>
									<Table.Body>
										<Table.Row class="*:whitespace-nowrap">
											<Table.Head>Name</Table.Head>
											<Table.Cell class="flex justify-end">{persistentVolumeClaim.name}</Table.Cell>
										</Table.Row>
										<Table.Row class="*:whitespace-nowrap">
											<Table.Head>Status</Table.Head>
											<Table.Cell class="flex justify-end">
												<Badge variant="outline">{persistentVolumeClaim.status}</Badge>
											</Table.Cell>
										</Table.Row>
										<Table.Row class="*:whitespace-nowrap">
											<Table.Head>Capacity</Table.Head>
											<Table.Cell class="flex justify-end"
												>{persistentVolumeClaim.capacity}</Table.Cell
											>
										</Table.Row>
										<Table.Row class="*:whitespace-nowrap">
											<Table.Head>Access Modes</Table.Head>
											<Table.Cell class="flex justify-end">
												{#each persistentVolumeClaim.accessModes as mode}
													<Badge variant="outline">{mode}</Badge>
												{/each}
											</Table.Cell>
										</Table.Row>
									</Table.Body>
								</Table.Root>
							</fieldset>
							{#if persistentVolumeClaim.storageClass}
								<fieldset>
									<legend class="text-center text-xs font-light">Storage Class</legend>
									<Table.Root>
										<Table.Body>
											<Table.Row class="*:whitespace-nowrap">
												<Table.Head>Name</Table.Head>
												<Table.Cell class="flex justify-end"
													>{persistentVolumeClaim.storageClass.name}</Table.Cell
												>
											</Table.Row>
											<Table.Row class="*:whitespace-nowrap">
												<Table.Head>Provisioner</Table.Head>
												<Table.Cell class="flex justify-end"
													>{persistentVolumeClaim.storageClass.provisioner}</Table.Cell
												>
											</Table.Row>
											<Table.Row class="*:whitespace-nowrap">
												<Table.Head>Reclaim Policy</Table.Head>
												<Table.Cell class="flex justify-end">
													<Badge variant="outline"
														>{persistentVolumeClaim.storageClass.reclaimPolicy}</Badge
													>
												</Table.Cell>
											</Table.Row>
											<Table.Row class="*:whitespace-nowrap">
												<Table.Head>Volume Binding Mode</Table.Head>
												<Table.Cell class="flex justify-end">
													<Badge variant="outline"
														>{persistentVolumeClaim.storageClass.volumeBindingMode}</Badge
													>
												</Table.Cell>
											</Table.Row>
										</Table.Body>
									</Table.Root>
								</fieldset>
							{/if}
						</HoverCard.Content>
					</HoverCard.Root>
					<Badge variant="outline">{persistentVolumeClaim.name}</Badge>
				</span>
			{/each}
		</Card.Footer>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Health</Card.Title>
		</Card.Header>
		<Card.Content class="text-3xl">
			{Math.round(health)}%
			<p class="text-xs text-muted-foreground">
				{numberOfHealthyPods} healthy pods over {numberOfPods} pods
			</p>
		</Card.Content>
		<Card.Footer>
			<Progress
				value={health}
				max={100}
				class={`${
					health > 62
						? 'bg-green-50 *:bg-green-700'
						: health > 38
							? 'bg-yellow-50 *:bg-yellow-500'
							: 'bg-red-50 *:bg-red-700'
				}`}
			/>
		</Card.Footer>
	</Card.Root>
{/snippet}

<div class="grid gap-3 p-3">
	{#if application.healthies !== application.pods.length}
		<Alert.Root variant="destructive">
			<span class="flex items-center gap-3">
				<Icon icon="radix-icons:exclamation-triangle" class="size-5" />

				<span>
					<Alert.Title>WARNING</Alert.Title>
					<Alert.Description>
						{application.pods.length - application.healthies} out of {application.pods.length} pods are
						unhealthy.
					</Alert.Description>
				</span>
			</span>
		</Alert.Root>
	{/if}
	{@render Identifier()}
	<div class="grid grid-cols-4 gap-3">
		{@render StatisticApplication(application)}
	</div>

	<div class="grid gap-3 [&>fieldset]:rounded-lg [&>fieldset]:border [&>fieldset]:p-3">
		{#if application.pods.length > 0}
			<fieldset>
				<legend>Pods</legend>
				<Table.Root>
					<Table.Header>
						<Table.Row class="*:text-xs *:font-light">
							<Table.Head>NAME</Table.Head>
							<Table.Head>PHASE</Table.Head>
							<Table.Head>READY</Table.Head>
							<Table.Head>RESTARTS</Table.Head>
							<Table.Head>LAST CONDITION</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each application.pods as pod}
							<Table.Row class="*:text-sm">
								<Table.Cell>{pod.name}</Table.Cell>
								<Table.Cell>
									<Badge variant="outline">{pod.phase}</Badge>
								</Table.Cell>
								<Table.Cell>
									{pod.ready}
								</Table.Cell>
								<Table.Cell>{pod.restarts}</Table.Cell>
								<Table.Cell>
									{#if pod.lastCondition}
										<div class="flex items-center gap-1">
											<HoverCard.Root openDelay={13}>
												<HoverCard.Trigger>
													<Icon icon="ph:info" class="size-4 text-blue-800" />
												</HoverCard.Trigger>
												<HoverCard.Content class=" min-w-96">
													<Table.Root>
														<Table.Body>
															<Table.Row>
																<Table.Head>Type</Table.Head>
																<Table.Cell>
																	<Badge variant="outline">{pod.lastCondition.type}</Badge>
																</Table.Cell>
															</Table.Row>
															<Table.Row>
																<Table.Head>Status</Table.Head>
																<Table.Cell>{pod.lastCondition.status}</Table.Cell>
															</Table.Row>
															<Table.Row>
																<Table.Head>Reason</Table.Head>
																<Table.Cell>{pod.lastCondition.reason}</Table.Cell>
															</Table.Row>
															<Table.Row>
																<Table.Head>Message</Table.Head>
																<Table.Cell>{pod.lastCondition.message}</Table.Cell>
															</Table.Row>
														</Table.Body>
													</Table.Root>
												</HoverCard.Content>
											</HoverCard.Root>
											<p>{pod.lastCondition.type}: {pod.lastCondition.status}</p>
										</div>
									{/if}
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</fieldset>
		{/if}
		{#if application.services?.length > 0}
			<fieldset>
				<legend>Services</legend>
				<Table.Root>
					<Table.Header>
						<Table.Row class="*:text-xs *:font-light">
							<Table.Head>NAME</Table.Head>
							<Table.Head>TYPE</Table.Head>
							<Table.Head>CLUSTER IP</Table.Head>
							<Table.Head>PORTS</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each application.services as service}
							<Table.Row class="*:text-sm">
								<Table.Cell>{service.name}</Table.Cell>
								<Table.Cell
									><Badge variant="outline" class="w-fit">{service.type}</Badge></Table.Cell
								>
								<Table.Cell>
									{service.clusterIp}
								</Table.Cell>
								<Table.Cell>
									<div class="flex flex-col">
										{#each service.ports as port}
											<span class="flex items-center gap-1">
												<HoverCard.Root openDelay={13}>
													<HoverCard.Trigger>
														<Icon icon="ph:info" class="size-4 text-blue-800" />
													</HoverCard.Trigger>
													<HoverCard.Content>
														<Table.Root>
															<Table.Body>
																<Table.Row>
																	<Table.Head>Protocol</Table.Head>
																	<Table.Cell>{port.protocol}</Table.Cell>
																</Table.Row>
																<Table.Row>
																	<Table.Head>Port</Table.Head>
																	<Table.Cell>{port.port}</Table.Cell>
																</Table.Row>
																<Table.Row>
																	<Table.Head>Target Port</Table.Head>
																	<Table.Cell>{port.targetPort}</Table.Cell>
																</Table.Row>
																<Table.Row>
																	<Table.Head>Node Port</Table.Head>
																	<Table.Cell>{port.nodePort}</Table.Cell>
																</Table.Row>
															</Table.Body>
														</Table.Root>
													</HoverCard.Content>
												</HoverCard.Root>

												<Badge variant="outline" class="w-fit">{port.protocol} {port.port}</Badge>
											</span>
										{/each}
									</div>
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</fieldset>
		{/if}
	</div>
</div>

{#snippet Identifier()}
	<span class="flex gap-3">
		<Icon icon="logos:kubernetes" class="h-full w-48" />
		<div class="flex flex-col justify-between">
			<div>
				<p class="text-xl">{application.namespace}/{application.name}</p>
				<Badge variant="outline">{application.type}</Badge>
			</div>
			<span class="flex flex-wrap gap-1">
				{#each Object.entries(application.labels) as [key, value]}
					<Badge
						variant="secondary"
						class="h-fit w-fit text-ellipsis whitespace-nowrap text-xs tracking-tight"
					>
						{key}: {value}
					</Badge>
				{/each}
			</span>
		</div>
	</span>
{/snippet}
