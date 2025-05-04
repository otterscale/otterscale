<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Card from '$lib/components/ui/card/index.js';
	import CreateKubernetes from '$lib/components/otterscale/management/general/create-kubernetes.svelte';
</script>

<!-- 
//TODO: add hover
Kubernetes is required for container orchestration and cluster management. Install it to enable container workload management.
-->

<Card.Root class="h-min w-[350px]">
	<Card.Header class="flex flex-col items-center justify-center">
		<div class="flex items-center gap-2">
			<Icon icon="simple-icons:kubernetes" class="h-24 w-24 p-2" style="color: #326de6" />
			<Icon icon="ph:plus-light" class="h-8 w-8" />
			<Icon icon="simple-icons:ceph" class="h-24 w-24 p-2" style="color: #f0424d" />
		</div>
		<Card.Title>Kubernetes and Ceph</Card.Title>
		<Card.Description>Deploy Container Platform and Storage</Card.Description>
	</Card.Header>
	<Card.Footer class="flex items-center justify-center">
		<CreateKubernetes />
	</Card.Footer>
</Card.Root>
