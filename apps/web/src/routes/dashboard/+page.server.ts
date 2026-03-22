import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

const API_URL = process.env.API_URL || 'http://localhost:3001';

interface Server {
	id: number;
	name: string;
	type: string;
	version: string;
	min_mem: number;
	max_mem: number;
	port: number;
}

export const load: PageServerLoad = async ({ locals, cookies }) => {
	const sessionCookie = cookies.get('session_id');
	let servers: Server[] = [];

	if (sessionCookie) {
		try {
			const res = await fetch(`${API_URL}/api/v1/servers`, {
				headers: {
					Cookie: `session_id=${sessionCookie}`
				}
			});
			if (res.ok) {
				servers = await res.json();
			}
		} catch {
			// Ignore errors, servers will be empty
		}
	}

	return {
		user: locals.user,
		servers
	};
};

export const actions: Actions = {
	logout: async ({ cookies }) => {
		const sessionCookie = cookies.get('session_id');

		if (sessionCookie) {
			try {
				// Call backend logout
				await fetch(`${API_URL}/api/v1/auth/logout`, {
					method: 'POST',
					headers: {
						Cookie: `session_id=${sessionCookie}`
					}
				});
			} catch {
				// Continue with cookie deletion even if backend fails
			}

			// Clear session cookie
			cookies.delete('session_id', { path: '/' });
		}

		throw redirect(302, '/');
	}
};
