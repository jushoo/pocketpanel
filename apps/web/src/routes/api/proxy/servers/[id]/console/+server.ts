import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const API_URL = process.env.API_URL || 'http://localhost:3001';

export const GET: RequestHandler = async ({ params, cookies, url }) => {
	const sessionCookie = cookies.get('session_id');
	const serverId = params.id;
	const lines = url.searchParams.get('lines') || '100';

	if (!sessionCookie) {
		return json({ error: 'Authentication required' }, { status: 401 });
	}

	try {
		const response = await fetch(`${API_URL}/api/v1/servers/${serverId}/console?lines=${lines}`, {
			headers: {
				Cookie: `session_id=${sessionCookie}`
			}
		});

		if (!response.ok) {
			const data = await response.json();
			return json({ error: data.error || 'Failed to fetch console' }, { status: response.status });
		}

		return json(await response.json());
	} catch (err) {
		return json({ error: 'An unexpected error occurred' }, { status: 500 });
	}
};
