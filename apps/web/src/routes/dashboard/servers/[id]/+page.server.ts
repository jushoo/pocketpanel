import { error, fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

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

export interface ServerWithStatus extends Server {
	running: boolean;
	pid?: number;
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

		const data = await response.json();
		return { 
			server: data.server as Server,
			running: data.running as boolean,
			pid: data.pid as number | undefined
		};
	} catch (err) {
		// Re-throw SvelteKit errors
		if (err && typeof err === 'object' && 'status' in err) {
			throw err;
		}
		throw error(500, 'An unexpected error occurred');
	}
};

export const actions: Actions = {
	start: async ({ params, cookies }) => {
		const sessionCookie = cookies.get('session_id');
		const serverId = params.id;

		if (!sessionCookie) {
			return fail(401, { error: 'Authentication required' });
		}

		try {
			const response = await fetch(`${API_URL}/api/v1/servers/${serverId}/start`, {
				method: 'POST',
				headers: {
					Cookie: `session_id=${sessionCookie}`
				}
			});

			if (!response.ok) {
				const data = await response.json();
				return fail(response.status, { error: data.error || 'Failed to start server' });
			}

			return { success: true, message: 'Server started successfully' };
		} catch (err) {
			return fail(500, { error: 'An unexpected error occurred' });
		}
	},

	stop: async ({ params, cookies }) => {
		const sessionCookie = cookies.get('session_id');
		const serverId = params.id;

		if (!sessionCookie) {
			return fail(401, { error: 'Authentication required' });
		}

		try {
			const response = await fetch(`${API_URL}/api/v1/servers/${serverId}/stop`, {
				method: 'POST',
				headers: {
					Cookie: `session_id=${sessionCookie}`,
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ force: false })
			});

			if (!response.ok) {
				const data = await response.json();
				return fail(response.status, { error: data.error || 'Failed to stop server' });
			}

			return { success: true, message: 'Server stopped successfully' };
		} catch (err) {
			return fail(500, { error: 'An unexpected error occurred' });
		}
	}
};
