<script lang="ts" module>
	import { page } from '$app/state';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import type { Application } from '../types';
	import Restart from './action-restart.svelte';
	import Scale from './action-scale.svelte';
</script>

<script lang="ts">
	let { application, scope, reloadManager }: { application: Application; scope: string; reloadManager: ReloadManager } = $props();
</script>

{#if page.data['feature-states.app-container']}
	<Layout.Actions>
		<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
		<Layout.ActionSeparator />
		<Layout.ActionItem>
			<Restart {application} {scope} {reloadManager} />
		</Layout.ActionItem>
		<Layout.ActionItem>
			<Scale {application} {scope} {reloadManager} />
		</Layout.ActionItem>
	</Layout.Actions>
{/if}
