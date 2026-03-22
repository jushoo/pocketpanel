import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

const API_URL = process.env.API_URL || 'http://localhost:3001';

export interface Server {
	id: number;
	name: string;
	type: string;
	version: string;
	min_mem: number;
	max_mem: number;
	port: number;
	created_at: string;
	updated_at: string;
}

export const load: PageServerLoad = async ({ params, cookies }) => {
	const sessionCookie = cookies.get('session_id');
	const serverId = params.id;

	if (!sessionCookie) {
		throw error(401, 'Authentication required');
	}

	try {
		const response = await fetch(`${API_URL}/api/v1/servers/${serverId}`, {
			headers: {
				Cookie: `session_id=${sessionCookie}`
			}
		});

		if (response.status === 404) {
			throw error(404, 'Server not found');
		}

		if (!response.ok) {
			throw error(500, 'Failed to load server');
		}

		const server: Server = await response.json();
		return { server };
	} catch (err) {
		// Re-throw SvelteKit errors
		if (err && typeof err === 'object' && 'status' in err) {
			throw err;
		}
		throw error(500, 'An unexpected error occurred');
	}
};
