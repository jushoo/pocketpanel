import { fail, redirect, isRedirect } from '@sveltejs/kit';
import type { Actions } from './$types';

const API_URL = process.env.API_URL || 'http://localhost:3001';

export const actions: Actions = {
	login: async ({ request, cookies }) => {
		const data = await request.formData();
		const username = data.get('username') as string;
		const password = data.get('password') as string;

		if (!username || !password) {
			return fail(400, {
				error: 'Please enter both username and password'
			});
		}

		try {
			const response = await fetch(`${API_URL}/api/v1/auth/login`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ username, password })
			});

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ error: 'Login failed' }));
				return fail(response.status, {
					error: errorData.error || 'Invalid credentials'
				});
			}

			// Extract session cookie from response
			const setCookieHeader = response.headers.get('set-cookie');
			if (setCookieHeader) {
				// Parse and forward the session cookie
				const sessionMatch = setCookieHeader.match(/session_id=([^;]+)/);
				if (sessionMatch) {
					const sessionId = sessionMatch[1];
					const maxAgeMatch = setCookieHeader.match(/Max-Age=(\d+)/);
					const maxAge = maxAgeMatch ? parseInt(maxAgeMatch[1]) : 24 * 60 * 60;
					
					cookies.set('session_id', sessionId, {
						path: '/',
						httpOnly: true,
						secure: process.env.NODE_ENV === 'production',
						maxAge: maxAge,
						sameSite: 'lax'
					});
				}
			}

			// Successful login - redirect to dashboard
			throw redirect(302, '/dashboard');
		} catch (error) {
			if (isRedirect(error)) {
				throw error;
			}
			
			console.error('Login error:', error);
			return fail(500, {
				error: 'An error occurred during login'
			});
		}
	}
};
