<script lang="ts" module>
	import type { Writable } from 'svelte/store';

	import type { Application_Release } from '$lib/api/application/v1/application_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	import Delete from './chart-action-delete-release.svelte';
	import Edit from './chart-action-edit-release.svelte';
	import Rollback from './chart-action-rollback-release.svelte';
</script>

<script lang="ts">
	let {
		release,
		releases = $bindable()
	}: {
		release: Application_Release;
		releases: Writable<Application_Release[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>
		{m.actions()}
	</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Edit {release} bind:releases />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Rollback {release} bind:releases />
	</Layout.ActionItem>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Delete {release} bind:releases />
	</Layout.ActionItem>
</Layout.Actions>
