<script lang="ts" module>
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Tabs from '$lib/components/custom/tabs';
	import { PrometheusDriver } from 'prometheus-query';
	import { Application } from './application';
	import { Storage } from './storage';
	import { System } from './overall';
</script>

<script lang="ts">
	let { client, scopes }: { client: PrometheusDriver; scopes: Scope[] } = $props();
</script>

<main class="no-user-select grid gap-4 p-4">
	<Tabs.Root value="storage">
		<Tabs.List>
			<Tabs.Trigger class="text-xl" value="overall">Overall</Tabs.Trigger>
			<Tabs.Trigger class="text-xl" value="storage">Storage</Tabs.Trigger>
			<Tabs.Trigger class="text-xl" value="application">Application</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="overall">
			<System {client} {scopes} />
		</Tabs.Content>
		<Tabs.Content value="storage">
			<Storage {client} {scopes} />
		</Tabs.Content>
		<Tabs.Content value="application">
			<Application {client} {scopes} />
		</Tabs.Content>
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
