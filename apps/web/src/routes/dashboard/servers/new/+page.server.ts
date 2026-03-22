import { fail, redirect, isRedirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

const API_URL = process.env.API_URL || 'http://localhost:3001';

export interface ServerType {
	id: string;
	name: string;
	description: string;
}

export const load: PageServerLoad = async ({ locals }) => {
	const serverTypes: ServerType[] = [
		{
			id: 'vanilla',
			name: 'Vanilla',
			description: 'Official Minecraft server'
		},
		{
			id: 'fabric',
			name: 'Fabric',
			description: 'Lightweight mod loader for Minecraft'
		}
	];

	return {
		user: locals.user,
		serverTypes
	};
};

export const actions: Actions = {
	create: async ({ request, cookies }) => {
		const sessionCookie = cookies.get('session_id');

		if (!sessionCookie) {
			return fail(401, {
				error: 'Authentication required',
				values: {}
			});
		}

		const data = await request.formData();

		const name = data.get('name') as string;
		const serverType = data.get('serverType') as string;
		const version = data.get('version') as string;
		const port = data.get('port') as string;
		const minMemory = data.get('minMemory') as string;
		const maxMemory = data.get('maxMemory') as string;

		// Validation
		if (!name || name.trim().length === 0) {
			return fail(400, {
				error: 'Server name is required',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}

		if (name.length > 100) {
			return fail(400, {
				error: 'Server name must be 100 characters or less',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}

		if (!serverType) {
			return fail(400, {
				error: 'Server type is required',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}

		if (!version) {
			return fail(400, {
				error: 'Server version is required',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}

		// Validate port if provided
		if (port) {
			const portNum = parseInt(port, 10);
			if (isNaN(portNum) || portNum < 1024 || portNum > 65535) {
				return fail(400, {
					error: 'Port must be between 1024 and 65535',
					values: { name, serverType, version, port, minMemory, maxMemory }
				});
			}
		}

		// Validate memory
		const minMemNum = parseInt(minMemory || '2', 10);
		const maxMemNum = parseInt(maxMemory || '4', 10);

		if (isNaN(minMemNum) || minMemNum < 1) {
			return fail(400, {
				error: 'Minimum memory must be at least 1 GB',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}

		if (isNaN(maxMemNum) || maxMemNum < minMemNum) {
			return fail(400, {
				error: 'Maximum memory must be greater than or equal to minimum memory',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}

		if (maxMemNum > 128) {
			return fail(400, {
				error: 'Maximum memory cannot exceed 128 GB',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}

		const serverPort = port ? parseInt(port, 10) : 25565;

		try {
			const response = await fetch(`${API_URL}/api/v1/servers`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Cookie: `session_id=${sessionCookie}`
				},
				body: JSON.stringify({
					name: name.trim(),
					type: serverType,
					version,
					port: serverPort,
					min_mem: minMemNum,
					max_mem: maxMemNum
				})
			});

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ error: 'Failed to create server' }));
				return fail(response.status, {
					error: errorData.error || 'Failed to create server',
					values: { name, serverType, version, port, minMemory, maxMemory }
				});
			}

			throw redirect(302, '/dashboard');
		} catch (error) {
			if (isRedirect(error)) {
				throw error;
			}

			console.error('Server creation error:', error);
			return fail(500, {
				error: 'An error occurred while creating the server',
				values: { name, serverType, version, port, minMemory, maxMemory }
			});
		}
	}
};
