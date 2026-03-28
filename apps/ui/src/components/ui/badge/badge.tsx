import { splitProps } from 'solid-js';
import type { JSX, Component } from 'solid-js';
import { cn } from '~/lib/utils';
import { cva, type VariantProps } from 'class-variance-authority';

const badgeVariants = cva(
	'inline-flex items-center rounded-md border px-2 py-0.5 text-xs font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
	{
		variants: {
			variant: {
				default: 'border-transparent bg-primary text-primary-foreground shadow hover:bg-primary/80',
				secondary: 'border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80',
				destructive:
					'border-transparent bg-destructive text-destructive-foreground shadow hover:bg-destructive/80',
				outline: 'text-foreground'
			}
		},
		defaultVariants: {
			variant: 'default'
		}
	}
);

export type BadgeVariant = VariantProps<typeof badgeVariants>['variant'];

export interface BadgeProps extends JSX.HTMLAttributes<HTMLSpanElement> {
	variant?: BadgeVariant;
}

export const Badge: Component<BadgeProps> = (props) => {
	const [local, rest] = splitProps(props, ['class', 'variant']);

	return (
		<span class={cn(badgeVariants({ variant: local.variant }), local.class)} {...rest} />
	);
};

export { badgeVariants };
