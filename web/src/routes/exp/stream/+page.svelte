<script lang="ts">
	import { getContext } from 'svelte';
	import { Code, ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { EnvironmentService } from '$gen/api/environment/v1/environment_pb';
	import { Button } from '$lib/components/ui/button';

	const transport: Transport = getContext('transport');
	const client = createClient(EnvironmentService, transport);

	let phase = $state('');
	let message = $state('');

	async function watchStatuses() {
		try {
			for await (const response of client.watchStatuses({})) {
				phase = response.phase;
				message = response.message;
			}
		} catch (error) {
			const connectError = error as ConnectError;
			if (connectError.code == Code.NotFound) {
				// 安裝包還沒送狀態
				// 可以先 postman 打 update status 測試
			}
			console.error(connectError);
		}
	}
</script>

<Button onclick={watchStatuses}>Watch</Button>

{phase}
{message}
