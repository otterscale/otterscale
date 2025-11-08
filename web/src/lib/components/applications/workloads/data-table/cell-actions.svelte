<script lang="ts" module>
	import { page } from '$app/state';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	import type { Application } from '../types';
	import Restart from './action-restart.svelte';
	import Scale from './action-scale.svelte';
</script>

<script lang="ts">
	let { application }: { application: Application } = $props();
</script>

{#if page.data['feature-states.app-container']}
	<Layout.Actions>
		<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
		<Layout.ActionSeparator />
		<Layout.ActionItem>
			<Restart {application} />
		</Layout.ActionItem>
		<Layout.ActionItem>
			<Scale {application} />
		</Layout.ActionItem>
	</Layout.Actions>
{/if}
