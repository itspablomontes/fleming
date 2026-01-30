import type { JSX } from "react";

import { Label } from "@/components/ui/label";

import type { ConsentForm } from "../consent-form-types";

interface ReasonStepProps {
	form: ConsentForm;
}

export function ReasonStep({ form }: ReasonStepProps): JSX.Element {
	return (
		<form.Field name="reason">
			{(field) => (
				<div className="space-y-3">
					<div className="space-y-2">
						<Label htmlFor="reason">Reason (optional)</Label>
						<p className="text-xs text-muted-foreground">
							You can skip this step.
						</p>
						<textarea
							id="reason"
							value={field.state.value}
							onChange={(event) => field.handleChange(event.target.value)}
							onBlur={field.handleBlur}
							placeholder="Explain why you need access"
							className="min-h-[110px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-500 focus-visible:ring-offset-2 dark:focus-visible:ring-cyan-400"
						/>
					</div>
					<p className="text-xs text-muted-foreground">
						A short reason helps patients approve requests faster.
					</p>
				</div>
			)}
		</form.Field>
	);
}
