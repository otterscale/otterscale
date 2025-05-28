<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import * as Tabs from '$lib/components/custom/tabs';

	import { System } from './system';
	import { Storage } from './storage';

	import type { Scope } from '$gen/api/scope/v1/scope_pb';

	let {
		client,
		scopes,
		instances
	}: { client: PrometheusDriver; scopes: Scope[]; instances: string[] } = $props();
</script>

<main class="no-user-select grid gap-4 p-4">
	<Tabs.Root value="storage">
		<Tabs.List>
			<Tabs.Trigger value="system">System</Tabs.Trigger>
			<Tabs.Trigger value="storage">Storage</Tabs.Trigger>
			<Tabs.Trigger value="application">Application</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="system">
			<System {client} {scopes} {instances} />
		</Tabs.Content>
		<Tabs.Content value="storage">
			<Storage {client} {scopes} {instances} />
		</Tabs.Content>
		<Tabs.Content value="application"></Tabs.Content>
	</Tabs.Root>
</main>

<style>
	.no-user-select {
		user-select: none;
		-webkit-user-select: none;
		-moz-user-select: none;
		-ms-user-select: none;
	}
</style>
