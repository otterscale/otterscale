<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';

	import { browser } from '$app/environment';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import type { Application, Application_Pod } from '$lib/api/application/v1/application_pb';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Empty } from '$lib/components/custom/table';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { m } from '$lib/paraglide/messages';

	import Actions from './cell-actions.svelte';
</script>

<script lang="ts">
	let {
		application,
		scope,
		namespace,
		reloadManager
	}: {
		application: Writable<Application>;
		scope: string;
		namespace: string;
		reloadManager: ReloadManager;
	} = $props();

	function openTerminal(pod: Application_Pod) {
		if (!browser) {
			return;
		}

		const searchParams = new URLSearchParams({
			scope,
			namespace: page.params.namespace ?? '',
			pod: pod.name,
			container: '',
			command: '/bin/sh'
		});

		const terminalUrl = `${resolve('/tty')}?${searchParams.toString()}`;
		const windowName = m.tty();

		const features = [
			'width=800',
			'height=600',
			'toolbar=no',
			'location=no',
			'menubar=no',
			'status=no',
			'scrollbars=no',
			'resizable=yes'
		].join(',');

		const newWindow = window.open(terminalUrl, windowName, features);

		if (newWindow) {
			newWindow.focus();
		}
	}
</script>

<Table.Root>
	<Table.Header>
		<Table.Row>
			<Table.Head>
				{m.name()}
			</Table.Head>
			<Table.Head>
				{m.phase()}
			</Table.Head>
			<Table.Head>
				{m.ready()}
			</Table.Head>
			<Table.Head>
				{m.restarts()}
			</Table.Head>
			<Table.Head>
				{m.conditions()}
			</Table.Head>
			<Table.Head>
				{m.terminal()}
			</Table.Head>
			<Table.Head></Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		{#each $application.pods as pod, index (index)}
			<Table.Row>
				<Table.Cell>{pod.name}</Table.Cell>
				<Table.Cell>
					<Badge variant="outline">{pod.phase}</Badge>
				</Table.Cell>
				<Table.Cell>
					{pod.ready}
				</Table.Cell>
				<Table.Cell>{pod.restarts}</Table.Cell>
				<Table.Cell>
					{#if pod.conditions}
						{@const trueConditions = pod.conditions.filter(
							(condition) => condition.status === 'True'
						)}
						<div class="flex flex-wrap gap-1">
							{#each trueConditions as trueCondition, index (index)}
								<Tooltip.Provider>
									<Tooltip.Root>
										<Tooltip.Trigger>
											<Badge
												variant={['Failed', 'FailureTarget'].includes(trueCondition.type)
													? 'destructive'
													: 'outline'}
											>
												{trueCondition.type}
											</Badge>
										</Tooltip.Trigger>
										<Tooltip.Content>
											{#if trueCondition.message}
												{trueCondition.message}
											{:else}
												{trueCondition.type}
											{/if}
										</Tooltip.Content>
									</Tooltip.Root>
								</Tooltip.Provider>
							{/each}
						</div>
					{/if}
				</Table.Cell>
				<Table.Cell>
					<Button variant="secondary" size="icon" onclick={() => openTerminal(pod)}>
						<Icon icon="ph:terminal-window" />
					</Button>
				</Table.Cell>
				<Table.Cell>
					<Actions {pod} {scope} {namespace} {reloadManager} />
				</Table.Cell>
			</Table.Row>
		{/each}
		{#if $application.pods.length === 0}
			<Table.Row>
				<Table.Cell colspan={6}>
					<Empty />
				</Table.Cell>
			</Table.Row>
		{/if}
	</Table.Body>
</Table.Root>
