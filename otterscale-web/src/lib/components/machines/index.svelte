<script lang="ts" module>
	import { activeScope } from '$lib/stores';
	import { Metals } from './metal/index';
</script>

{#if $activeScope}
	{#key $activeScope.uuid}
		<Metals />
	{/key}
{/if}
