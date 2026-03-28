import { createEffect, createSignal } from 'solid-js';
import { Sun, Moon } from 'lucide-solid';
import { Button } from '~/components/ui/button';

export function ModeToggle() {
	const [theme, setTheme] = createSignal<'light' | 'dark'>('light');

	createEffect(() => {
		const savedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null;
		const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
		const initialTheme = savedTheme || (prefersDark ? 'dark' : 'light');
		setTheme(initialTheme);
		updateTheme(initialTheme);
	});

	const updateTheme = (newTheme: 'light' | 'dark') => {
		const root = document.documentElement;
		if (newTheme === 'dark') {
			root.classList.add('dark');
		} else {
			root.classList.remove('dark');
		}
	};

	const toggleMode = () => {
		const newTheme = theme() === 'light' ? 'dark' : 'light';
		setTheme(newTheme);
		localStorage.setItem('theme', newTheme);
		updateTheme(newTheme);
	};

	return (
		<div class="fixed top-4 right-4 z-50">
			<Button onClick={toggleMode} variant="outline" size="icon">
				<Sun
					class={`h-[1.2rem] w-[1.2rem] scale-100 rotate-0 transition-all ${theme() === 'dark' ? 'scale-0 -rotate-90' : ''}`}
				/>
				<Moon
					class={`absolute h-[1.2rem] w-[1.2rem] scale-0 rotate-90 transition-all ${theme() === 'dark' ? 'scale-100 rotate-0' : ''}`}
				/>
				<span class="sr-only">Toggle theme</span>
			</Button>
		</div>
	);
}
