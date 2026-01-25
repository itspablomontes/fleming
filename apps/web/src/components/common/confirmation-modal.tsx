import { AlertTriangle, Info } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";

type ConfirmationVariant = "warning" | "info";

interface ConfirmationModalProps {
	isOpen: boolean;
	onOpenChange: (open: boolean) => void;
	onConfirm: () => void;
	title: string;
	message: string;
	confirmLabel?: string;
	cancelLabel?: string;
	variant?: ConfirmationVariant;
	isConfirming?: boolean;
}

const variantStyles: Record<ConfirmationVariant, { icon: typeof AlertTriangle; iconClass: string }> = {
	warning: {
		icon: AlertTriangle,
		iconClass: "text-amber-600 dark:text-amber-400",
	},
	info: {
		icon: Info,
		iconClass: "text-cyan-600 dark:text-cyan-400",
	},
};

export function ConfirmationModal({
	isOpen,
	onOpenChange,
	onConfirm,
	title,
	message,
	confirmLabel = "Confirm",
	cancelLabel = "Cancel",
	variant = "warning",
	isConfirming = false,
}: ConfirmationModalProps) {
	const Icon = variantStyles[variant].icon;

	return (
		<Dialog open={isOpen} onOpenChange={onOpenChange}>
			<DialogContent className="sm:max-w-[480px] bg-white dark:bg-gray-950 border-cyan-200 dark:border-cyan-900">
				<DialogHeader>
					<div className="mx-auto bg-cyan-100 dark:bg-cyan-900/30 p-3 rounded-full mb-4">
						<Icon className={cn("h-6 w-6", variantStyles[variant].iconClass)} />
					</div>
					<DialogTitle className="text-center text-xl text-cyan-900 dark:text-cyan-50">
						{title}
					</DialogTitle>
					<DialogDescription className="text-center text-gray-500 dark:text-gray-400">
						{message}
					</DialogDescription>
				</DialogHeader>
				<DialogFooter className="sm:justify-center">
					<Button
						variant="ghost"
						onClick={() => onOpenChange(false)}
						disabled={isConfirming}
					>
						{cancelLabel}
					</Button>
					<Button
						onClick={onConfirm}
						disabled={isConfirming}
						className="bg-cyan-600 hover:bg-cyan-700 text-white"
					>
						{isConfirming ? "Please wait..." : confirmLabel}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
