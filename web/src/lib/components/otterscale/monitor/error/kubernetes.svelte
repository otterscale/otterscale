<script lang="ts">
	import Icon from '@iconify/svelte';
	import CreateKubernetes from '$lib/components/otterscale/management/general/create-kubernetes.svelte';
</script>

<div
	class="flex animate-pulse flex-col items-center justify-center p-24 hover:animate-none hover:cursor-pointer"
>
	<Icon icon="simple-icons:kubernetes" class="h-72 w-72" />
	<h1 class="text-lg">Kubernetes is not installed</h1>
	<p class="w-72 text-xs text-muted-foreground">
		Kubernetes is required for container orchestration and cluster management. Install it to enable
		container workload management.
	</p>
	<CreateKubernetes />
</div>
