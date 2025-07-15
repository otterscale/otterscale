<script lang="ts">
	import { LoginForm } from '$lib/components/login';

	const { data } = $props();
</script>

<div class="bg-muted flex min-h-svh flex-col items-center justify-center p-6 md:p-10">
	<div class="w-full max-w-sm md:max-w-3xl">
		<LoginForm {data} />
	</div>
</div>
