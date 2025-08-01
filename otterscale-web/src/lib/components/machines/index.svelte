<script lang="ts" module>
	import { activeScope } from '$lib/stores';
	import { Metals } from './metal/index';
</script>

{#key $activeScope.uuid}
	<Metals />
{/key}
