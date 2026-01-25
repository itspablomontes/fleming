import { Check } from "lucide-react";
import type { JSX } from "react";

import { cn } from "@/lib/utils";

interface StepIndicatorProps {
	steps: readonly string[];
	currentStep: number;
	className?: string;
}

export function StepIndicator({
	steps,
	currentStep,
	className,
}: StepIndicatorProps): JSX.Element {
	const progress =
		steps.length > 1 ? (currentStep / (steps.length - 1)) * 100 : 0;

	return (
		<div className={cn("space-y-3", className)}>
			<div className="flex items-center justify-between text-xs text-muted-foreground">
				<span>
					Step {currentStep + 1} of {steps.length}
				</span>
				<span className="sr-only" aria-live="polite">
					Current step: {steps[currentStep]}
				</span>
			</div>
			<div className="h-1.5 w-full rounded-full bg-muted">
				<div
					className="h-full rounded-full bg-primary transition-all duration-300"
					style={{ width: `${progress}%` }}
				/>
			</div>
			<div className="grid gap-3 sm:grid-cols-5">
				{steps.map((label, index) => {
					const isComplete = index < currentStep;
					const isActive = index === currentStep;

					return (
						<div key={label} className="flex items-center gap-2">
							<div
								className={cn(
									"flex h-7 w-7 items-center justify-center rounded-full border text-xs font-semibold",
									isComplete &&
										"border-primary bg-primary text-primary-foreground",
									isActive &&
										!isComplete &&
										"border-primary text-primary ring-2 ring-primary/30",
									!isComplete &&
										!isActive &&
										"border-border text-muted-foreground",
								)}
								aria-hidden="true"
							>
								{isComplete ? <Check className="h-3.5 w-3.5" /> : index + 1}
							</div>
							<span
								className={cn(
									"text-xs font-medium",
									isActive && "text-foreground",
									!isActive && "text-muted-foreground",
								)}
							>
								{label}
							</span>
						</div>
					);
				})}
			</div>
		</div>
	);
}
