import createClient from 'openapi-fetch';
import { vi } from 'vitest';
import type { paths } from '$lib/api/schema';

/** A JSON response, 200 unless `status` says otherwise. */
export const json = (body: unknown, status = 200) => Response.json(body, { status });

/**
 * An API client whose `fetch` answers with the given responses, in order (a promise answers
 * when it resolves). `request(n)` is the n-th request it was sent, from 0.
 */
export function stubApi(...responses: (Response | Promise<Response>)[]) {
	const fetch = vi.fn<(request: Request) => Promise<Response>>();
	for (const response of responses) fetch.mockReturnValueOnce(Promise.resolve(response));
	const client = createClient<paths>({ baseUrl: 'http://test/api/v1', fetch });
	const request = (n: number) => {
		const call = fetch.mock.calls[n];
		if (!call) throw new Error(`request ${n} was never sent (${fetch.mock.calls.length} were)`);
		return call[0];
	};
	return { client, fetch, request };
}
