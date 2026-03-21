import type { Handle } from '@sveltejs/kit';

const API_URL = process.env.API_URL || 'http://localhost:3001';

export const handle: Handle = async ({ event, resolve }) => {
	// Check if user is authenticated for protected routes
	if (event.url.pathname.startsWith('/dashboard')) {
		const sessionCookie = event.cookies.get('session_id');

		if (!sessionCookie) {
			return new Response(null, {
				status: 302,
				headers: { Location: '/' }
			});
		}

		// Verify session with backend
		try {
			const response = await fetch(`${API_URL}/api/v1/me`, {
				headers: {
					Cookie: `session_id=${sessionCookie}`
				}
			});

			if (!response.ok) {
				return new Response(null, {
					status: 302,
					headers: { Location: '/' }
				});
			}

			const user = await response.json();
			event.locals.user = user;
		} catch {
			return new Response(null, {
				status: 302,
				headers: { Location: '/' }
			});
		}
	}

	// Redirect authenticated users away from login page
	if (event.url.pathname === '/') {
		const sessionCookie = event.cookies.get('session_id');

		if (sessionCookie) {
			try {
				const response = await fetch(`${API_URL}/api/v1/me`, {
					headers: {
						Cookie: `session_id=${sessionCookie}`
					}
				});

				if (response.ok) {
					return new Response(null, {
						status: 302,
						headers: { Location: '/dashboard' }
					});
				}
			} catch {
				// Session invalid, continue to login
			}
		}
	}

	return resolve(event);
};
