import type { RequestEvent } from '@sveltejs/kit';

export async function POST(event: RequestEvent): Promise<Response> {
	try {
		const { serviceUri, modelName, modelIdentifier, prompt, max_tokens, temperature } =
			await event.request.json();

		const upstream = await fetch(`${serviceUri}/v1/completions`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				'OtterScale-Model-Name': modelName
			},
			body: JSON.stringify({
				model: modelIdentifier,
				prompt,
				max_tokens,
				temperature
			})
		});

		const body = await upstream.text();

		return new Response(body, {
			status: upstream.status,
			headers: upstream.headers
		});
	} catch (error) {
		return new Response(JSON.stringify({ error: error ?? 'Internal Server Error' }), {
			status: 500,
			headers: { 'Content-Type': 'application/json' }
		});
	}
}
