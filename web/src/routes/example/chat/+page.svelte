<script lang="ts">
	import Prompting from '$lib/components/prompting/index.svelte';

	const namespace = 'namespace';
	const modelName = 'meta-llama/Llama-3.1-8B-Instruct';
</script>

<Prompting {namespace} {modelName} />
