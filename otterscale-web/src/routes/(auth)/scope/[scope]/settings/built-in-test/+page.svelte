<script lang="ts">
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import * as Alert from '$lib/components/ui/alert';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
	import { BuiltInSelfTest } from '$lib/components/settings/built-in-self-test/index';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.settings(page.params.scope)],
		current: { title: m.built_in_test(), url: '' }
	});
</script>

{#if $activeScope}
	<BuiltInSelfTest />
{/if}
