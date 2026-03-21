import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

const API_URL = process.env.API_URL || 'http://localhost:3001';

export const load: PageServerLoad = async ({ locals }) => {
	return {
		user: locals.user
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
						'Cookie': `session_id=${sessionCookie}`
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
