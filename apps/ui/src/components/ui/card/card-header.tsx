import { splitProps } from 'solid-js';
import type { JSX, Component } from 'solid-js';
import { cn } from '~/lib/utils';

export interface CardHeaderProps extends JSX.HTMLAttributes<HTMLDivElement> {}

export const CardHeader: Component<CardHeaderProps> = (props) => {
	const [local, rest] = splitProps(props, ['class']);

	return (
		<div
			data-slot="card-header"
			class={cn('flex flex-col space-y-1.5 p-6', local.class)}
			{...rest}
		/>
	);
};
