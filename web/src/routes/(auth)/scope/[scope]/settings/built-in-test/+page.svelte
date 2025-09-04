<script lang="ts">
	import { page } from '$app/state';
	import { BuiltInSelfTest } from '$lib/components/settings/built-in-self-test/index';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.settings(page.params.scope)],
		current: { title: m.built_in_test(), url: '' },
	});
</script>

{#if $activeScope}
	<BuiltInSelfTest />
{/if}
