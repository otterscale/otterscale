<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';

	import Actions from './cell-actions.svelte';

	import { browser } from '$app/environment';
	import { page } from '$app/state';
	import type { Application, Application_Pod } from '$lib/api/application/v1/application_pb';
	import * as Table from '$lib/components/custom/table';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { m } from '$lib/paraglide/messages';
	import { staticPaths } from '$lib/path';
	import { currentKubernetes } from '$lib/stores';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		application,
	}: {
		application: Writable<Application>;
	} = $props();

	function openTerminal(pod: Application_Pod) {
		if (!browser) {
			return;
		}

		const searchParams = new URLSearchParams({
			scope: $currentKubernetes?.scope ?? '',
			facility: $currentKubernetes?.name ?? '',
			namespace: page.params.namespace ?? '',
			pod: pod.name,
			container: '',
			command: '/bin/sh',
		});

		const terminalUrl = staticPaths.tty.url + `?${searchParams.toString()}`;
		const windowName = m.tty();

		const features = [
			'width=800',
			'height=600',
			'toolbar=no',
			'location=no',
			'menubar=no',
			'status=no',
			'scrollbars=no',
			'resizable=yes',
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
				{m.last_condition()}
			</Table.Head>
			{#if page.data['feature-states-app-container']}
				<Table.Head>
					{m.terminal()}
				</Table.Head>
			{/if}
			<Table.Head></Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		{#each $application.pods as pod}
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
					{#if pod.lastCondition}
						{#if pod.lastCondition.reason || pod.lastCondition.message}
							<div class="text-destructive flex items-center gap-2">
								<Badge variant="destructive" class={pod.lastCondition.reason ? 'visible' : 'hidden'}>
									{pod.lastCondition.reason}
								</Badge>
								<Tooltip.Provider>
									<Tooltip.Root>
										<Tooltip.Trigger>
											<p
												class={cn(
													pod.lastCondition.message ? 'max-w-[1000px] truncate' : 'hidden',
												)}
											>
												{pod.lastCondition.message}
											</p>
										</Tooltip.Trigger>
										<Tooltip.Content class="max-w-[77vw] overflow-auto">
											{pod.lastCondition.message}
										</Tooltip.Content>
									</Tooltip.Root>
								</Tooltip.Provider>
							</div>
						{:else}
							<Badge variant="outline">
								{pod.lastCondition.type}
							</Badge>
						{/if}
					{/if}
				</Table.Cell>
				{#if page.data['feature-states-app-container']}
					<Table.Cell>
						<Button variant="secondary" size="icon" onclick={() => openTerminal(pod)}>
							<Icon icon="ph:terminal-window" />
						</Button>
					</Table.Cell>
				{/if}
				<Table.Cell class="p-0">
					<Actions {pod} namespace={$application.namespace} />
				</Table.Cell>
			</Table.Row>
		{/each}
		{#if $application.pods.length === 0}
			<Table.Row>
				<Table.Cell colspan={6}>
					<Table.Empty />
				</Table.Cell>
			</Table.Row>
		{/if}
	</Table.Body>
</Table.Root>
