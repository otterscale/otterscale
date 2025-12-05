import type { RequestEvent } from '@sveltejs/kit';

export async function POST(event: RequestEvent): Promise<Response> {
	const { serviceUri, modelName, modelIdentifier, prompt, max_tokens, temperature } =
		await event.request.json();

	const response = await fetch(`${serviceUri}/v1/completions`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', 'OtterScale-Model-Name': modelName },
		body: JSON.stringify({
			model: modelIdentifier,
			prompt: prompt,
			max_tokens: max_tokens,
			temperature: temperature
		})
	});

	const body = await response.text();

	return new Response(body, {
		status: response.status,
		headers: response.headers
	});
}
