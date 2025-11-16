<script lang="ts" module>
	import { page } from '$app/state';
	import type { Application_Pod } from '$lib/api/application/v1/application_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Log from './action-log.svelte';
</script>

<script lang="ts">
	let { pod, scope, namespace, reloadManager }: { pod: Application_Pod; scope: string; namespace: string; reloadManager: ReloadManager } = $props();
</script>

{#if page.data['feature-states.app-container']}
	<Layout.Actions>
		<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
		<Layout.ActionSeparator />
		<Layout.ActionItem>
			<Log {pod} {scope} {namespace}  />
		</Layout.ActionItem>
		<Layout.ActionItem>
			<Delete {pod} {scope} {namespace} {reloadManager} />
		</Layout.ActionItem>
	</Layout.Actions>
{/if}
