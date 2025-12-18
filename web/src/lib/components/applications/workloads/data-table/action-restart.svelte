<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import type { Application } from '../types';
</script>

<script lang="ts">
	let {
		application,
		scope,
		reloadManager,
		closeActions
	}: {
		application: Application;
		scope: string;
		reloadManager: ReloadManager;

		closeActions: () => void;
	} = $props();

	const transport: Transport = getContext('transport');

	const applicationClient = createClient(ApplicationService, transport);
	let loading = $state(false);

	async function restartApplication() {
		const request = {
			scope: scope,
			name: application.name,
			namespace: application.namespace,
			type: application.type
		};

		toast.promise(() => applicationClient.restartApplication(request), {
			loading: `Restarting application ${request.name}...`,
			success: () => {
				// Force reload to refresh data
				reloadManager.force();
				return `Successfully restarted application ${request.name}.`;
			},
			error: (e) => {
				const msg = `Failed to restart application ${request.name}.`;
				toast.error(msg, {
					description: (e as ConnectError).message.toString(),
					duration: Number.POSITIVE_INFINITY
				});
				return msg;
			}
		});
	}

	async function handleClick() {
		loading = true;

		try {
			await restartApplication();
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex items-center justify-end gap-1">
	<button
		onclick={async () => {
			await handleClick();
			closeActions();
		}}
		disabled={loading}
		class="flex items-center gap-1 disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50"
	>
		{#if loading}
			<Icon icon="ph:spinner-gap" class="animate-spin" />
			{m.please_wait()}
		{:else}
			<Icon icon="ph:arrow-clockwise" class="text-accent-foreground" />
			{m.restarts()}
		{/if}
	</button>
</div>
