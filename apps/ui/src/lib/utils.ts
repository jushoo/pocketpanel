import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function formatMemory(mb: number): string {
	if (mb >= 1024) {
		return `${(mb / 1024).toFixed(1)}GB`;
	}
	return `${mb}MB`;
}
