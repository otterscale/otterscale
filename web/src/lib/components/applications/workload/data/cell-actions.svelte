<script lang="ts" module>
	import Delete from './action-delete.svelte';
	import Log from './action-log.svelte';

	import { page } from '$app/state';
	import type { Application_Pod } from '$lib/api/application/v1/application_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { namespace, pod }: { namespace: string; pod: Application_Pod } = $props();
</script>

{#if page.data['feature-states.app-container']}
	<Layout.Actions>
		<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
		<Layout.ActionSeparator />
		<Layout.ActionItem>
			<Log {namespace} {pod} />
		</Layout.ActionItem>
		<Layout.ActionItem>
			<Delete {namespace} {pod} />
		</Layout.ActionItem>
	</Layout.Actions>
{/if}
